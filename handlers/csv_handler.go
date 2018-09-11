package handlers

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sscp/telemetry/log"

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
	csvWrite   *csv.Writer
	dmBuffer   *DataMapBuffer
	keyToIndex map[string]interface{}
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
	cw.createFile(ctx, runName, startTime)
	cw.setupWriter()
}

// createFile creates a .csv file to write to and log.Fatals if it errors (unlikely)
func (cw *CSVWriter) createFile(ctx context.Context, runName string, startTime time.Time) {
	filename := GetCSVFileName(runName, startTime)
	var err error
	cw.file, err = os.Create(filepath.Join(cw.folderPath, filename))
	if err != nil {
		log.Error(ctx, err, "Could not create csv file")
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
	cw.csvWrite = csv.NewWriter(cw.buffer)
	cw.dmBuffer = NewDataMessageBuffer(cw.writeData, 10)
}

// flushDataBuffer writes all the data in the DataMessage buffer to the
// buffered writer
func (cw *CSVWriter) writeData(ctx context.Context, data []map[string]interface{}) {
	if len(data) > 0 && (cw.keys == nil || len(cw.keys) == 0) {
		cw.keyToIndex = make(map[string]int, len(data[0]))
	}
	for _, dataMap := range data {
		csvLine := make([]string, len(dataMap))
		for key, value := range dataMap {
			index, ok := cw.keyToIndex[key]
			if !ok {
				cw.keyToIndex[key] = len(cw.keyToIndex)
			}
			csvLine[index] = value.String()
		}
		cw.csvWrite.Marshall(csvLine)
	}
}

// HandleData is called on every new DataMessage from the collector and adds
// the new DataMessage to the buffer of DataMessages and flushes the buffer if
// it is full
func (cw *CSVWriter) HandleData(ctx context.Context, data map[string]interface{}) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CSVWriter/HandleData")
	defer span.Finish()

	cw.dmBuffer.AddData(ctx, data)
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
	cw.dmBuffer.Flush(ctx)
	cw.buffer.Flush()

	csvLine := make([]string, len(cw.keyToIndex))
	for key, index := range keyToIndex {
		csvLine[index] = key
	}
	cw.csvWrite.Marshall(csvLine)

	cw.file.Close()
}