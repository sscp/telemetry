package collector

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	sundaeproto "github.com/sscp/telemetry/collector/sundae"

	"github.com/gocarina/gocsv"
	"github.com/opentracing/opentracing-go"
)

// GetCSVFileName wraps GetBaseFileName, in handlers.go, and adds the .csv to the end
func GetCSVFileName(runName string, startTime time.Time) string {
	return fmt.Sprintf("%v.csv", GetBaseFileName(runName, startTime))
}

// CSVConfig contains config info for the CSVWriter
type CSVConfig struct {
	Folder string
}

// CSVWriter is a DataHandler (handlers.go) that writes to CSV files
type CSVWriter struct {
	folderPath string
	file       *os.File
	buffer     *bufio.Writer
	csvWrite   *gocsv.SafeCSVWriter
	dmBuffer   *DataMessageBuffer
}

// dataBufferSize is the default size of the queues that lead to each handler
const dataBufferSize = 10

// NewCSVWriter returns an instantiated CSVWriter as a DataHandler interface
func NewCSVWriter(cfg CSVConfig) (DataHandler, error) {
	err := os.MkdirAll(cfg.Folder, os.ModePerm) // Create folder if it doesn't exist
	if err != nil {
		return nil, err
	}
	return &CSVWriter{folderPath: cfg.Folder}, nil
}

// HandleStartRun is called when collector starts recording a run and creates
// the CSV file and sets up all buffers
func (cw *CSVWriter) HandleStartRun(ctx context.Context, runName string, startTime time.Time) {
	cw.createFile(runName, startTime)
	cw.setupWriter()
}

// createFile creates a .csv file to write to and log.Fatals if it errors (unlikely)
func (cw *CSVWriter) createFile(runName string, startTime time.Time) {
	filename := GetCSVFileName(runName, startTime)
	var err error
	cw.file, err = os.Create(filepath.Join(cw.folderPath, filename))
	if err != nil {
		log.Fatal(err)
	}
}

// setupWriter creates the buffered writer for the raw file (reduces sys
// calls), and sets up the *DataMessage buffer, reduces calls to MarshalCSV,
// which is has the non-negligible cost of struct reflection to read the struct
// tags to determine how to convert each struct into a csv row. This also
// writes the first row of headers to the CSv file as every future call will
// simply append rows of data without headers.
func (cw *CSVWriter) setupWriter() {
	cw.buffer = bufio.NewWriter(cw.file)
	cw.csvWrite = gocsv.DefaultCSVWriter(cw.buffer)
	cw.dmBuffer = NewDataMessageBuffer(cw.writeData, 10)
	// Write only the headers because all future data will be written
	// without headers
	empData := []*sundaeproto.DataMessage{}
	gocsv.MarshalCSV(&empData, cw.csvWrite)
}

// flushDataBuffer writes all the data in the DataMessage buffer to the
// buffered writer
func (cw *CSVWriter) writeData(data []*sundaeproto.DataMessage) {
	// Write all data up until the current index. We can't write all the
	// data because there might be already written *DataMessages beyond the
	// current index
	gocsv.MarshalCSVWithoutHeaders(data, cw.csvWrite)
}

// HandleData is called on every new DataMessage from the collector and adds
// the new DataMessage to the buffer of DataMessages and flushes the buffer if
// it is full
func (cw *CSVWriter) HandleData(ctx context.Context, data *sundaeproto.DataMessage) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CSVWriter/HandleData")
	defer span.Finish()
	cw.dmBuffer.AddData(data)
}

// HandleDroppedData is called whenever CSVWriter falls behind and currently
// does nothing other than report a span to tracing
func (cw *CSVWriter) HandleDroppedData(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CSVWriter/HandleDroppedPacket")
	defer span.Finish()
}

// HandleEndRun is called by collector when data collection stops and the queue
// is empty and all data is processed. It flushes all the buffers to the .csv
// file. This must happen in the order of data buffer, buffered writer, file
// close to ensure no data loss.
func (cw *CSVWriter) HandleEndRun(ctx context.Context, endTime time.Time) {
	cw.dmBuffer.Flush()
	cw.buffer.Flush()
	cw.file.Close()
}
