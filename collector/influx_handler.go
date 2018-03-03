package collector

import (
	"context"
	"fmt"
	"log"
	"time"

	sundaeproto "github.com/sscp/telemetry/collector/sundae"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/opentracing/opentracing-go"

	// Fork of https://github.com/fatih/structs/ that adds an "indirect"
	// option to dereference pointers to get values, not pointers in map
	"github.com/jackbeasley/structs"
)

const databaseName = "sundae"

// queryDB convenience function to query influx
func queryDB(clnt influx.Client, cmd string) (res []influx.Result, err error) {
	q := influx.Query{
		Command:  cmd,
		Database: databaseName,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// InfluxConfig contains config info for the InfluxWriter
type InfluxConfig struct {
	Addr     string
	Username string
	Password string
}

// InfluxWriter is a DataHandler (handlers.go) that writes to InfluxDB
type InfluxWriter struct {
	config   InfluxConfig
	runName  string
	client   influx.Client
	dataBuf  []*sundaeproto.DataMessage
	bufIndex int
}

// NewInfluxWriter returns an instantiated InfluxWriter as a DataHandler interface
func NewInfluxWriter(cfg InfluxConfig) (DataHandler, error) {
	return &InfluxWriter{config: cfg}, nil
}

// HandleStartRun is called when collector starts recording a run and creates
// the CSV file and sets up all buffers
func (cw *InfluxWriter) HandleStartRun(ctx context.Context, runName string, startTime time.Time) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InfluxWriter/HandleStartRun")
	defer span.Finish()
	// Create a new HTTPClient
	var err error
	cw.client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     cw.config.Addr,
		Username: cw.config.Username,
		Password: cw.config.Password,
	})
	if err != nil {
		log.Printf("Could not create influx client: %v", err)
	}

	_, resp, err := cw.client.Ping(1 * time.Second)
	if err != nil {
		log.Printf("Could not ping influx: %v", err)
	}
	log.Printf("Connected to Influx at %v, version: %v", cw.config.Addr, resp)

	_, err = queryDB(cw.client, fmt.Sprintf("CREATE DATABASE %s", databaseName))
	if err != nil {
		log.Printf("Error creating database: %v", err)
	}

	cw.dataBuf = make([]*sundaeproto.DataMessage, 10)
	cw.bufIndex = 0
	cw.runName = runName
}

// setupWriter connects to the database and sets up the *DataMessage buffer,
// which creates larger batches of points to send to the database at once as a
// point batch
func (cw *InfluxWriter) setupWriter() error {
	return nil
}

// flushDataBuffer writes all the data in the DataMessage buffer to influx as a
// point batch
func (cw *InfluxWriter) flushDataBuffer() error {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  databaseName,
		Precision: "ns",
	})
	if err != nil {
		return err
	}
	for _, dm := range cw.dataBuf[0:cw.bufIndex] {
		// Convert struct to map[string]interface{}
		dataFields := structs.Map(dm)
		// Create a point and add to batch
		tags := map[string]string{"run_name": cw.runName}
		// TimeCollected is always set when deserialized by collector
		pt, err := influx.NewPoint("car_state", tags, dataFields, time.Unix(0, dm.GetTimeCollected()))
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}
	err = cw.client.Write(bp)
	if err != nil {
		return err
	}
	cw.bufIndex = 0
	return nil
}

// HandleData is called on every new DataMessage from the collector and adds
// the new DataMessage to the buffer of DataMessages and flushes the buffer if
// it is full
func (cw *InfluxWriter) HandleData(ctx context.Context, data *sundaeproto.DataMessage) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InfluxWriter/HandleData")
	defer span.Finish()
	cw.dataBuf[cw.bufIndex] = data
	cw.bufIndex++
	if cw.bufIndex >= len(cw.dataBuf) {
		cw.flushDataBuffer()
	}
}

// HandleDroppedData is called whenever InfluxWriter falls behind and currently
// does nothing other than report a span to tracing
func (cw *InfluxWriter) HandleDroppedData(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InfluxWriter/HandleDroppedPacket")
	defer span.Finish()
}

// HandleEndRun is called by collector when data collection stops and the queue
// is empty and all data is processed. It flushes all the buffers to the .csv
// file. This must happen in the order of data buffer, buffered writer, file
// close to ensure no data loss.
func (cw *InfluxWriter) HandleEndRun(ctx context.Context, endTime time.Time) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InfluxWriter/HandleEndRun")
	defer span.Finish()

	cw.flushDataBuffer()
	cw.client.Close()
	cw.runName = ""
}
