package handlers

import (
	"context"
	"fmt"
	"math"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/opentracing/opentracing-go"

	"github.com/sscp/telemetry/events"
	"github.com/sscp/telemetry/log"
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
	deBuffer *DataEventBuffer
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
		log.Error(ctx, err, "Could not create influx client")
		return
	}

	_, resp, err := cw.client.Ping(1 * time.Second)
	if err != nil {
		log.Error(ctx, err, "Could not ping influx")
		return
	}
	log.Event(ctx, fmt.Sprintf("Connected to Influx at %v, version: %v", cw.config.Addr, resp))

	_, err = queryDB(cw.client, fmt.Sprintf("CREATE DATABASE %s", databaseName))
	if err != nil {
		log.Error(ctx, err, "Error creating influx database")
		return
	}

	cw.deBuffer = NewDataEventBuffer(cw.writeData, 10)
	cw.runName = runName
}

// setupWriter connects to the database and sets up the *DataMessage buffer,
// which creates larger batches of points to send to the database at once as a
// point batch
func (cw *InfluxWriter) setupWriter() error {
	return nil
}

// writeData writes all the data in the DataMessage buffer to influx as a
// point batch
func (cw *InfluxWriter) writeData(ctx context.Context, data []events.DataEvent) {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  databaseName,
		Precision: "ns",
	})
	if err != nil {
		log.Error(ctx, err, "Error creating batch points")
		return
	}
	for _, dataEvent := range data {
		for key, value := range dataEvent.Data {
			if val, ok := value.(float32); ok {
				if math.IsNaN(float64(val)) {
					delete(dataEvent.Data, key)
				}
			}
			if val, ok := value.(float64); ok {
				if math.IsNaN(val) {
					delete(dataEvent.Data, key)
				}
			}

		}
		// Create a point and add to batch
		tags := map[string]string{"run_name": cw.runName}
		// TimeCollected is always set when deserialized by collector
		pt, err := influx.NewPoint("car_state", tags, dataEvent.Data, dataEvent.GetCollectedTime())
		if err != nil {
			log.Error(ctx, err, "Error creating influx point")
			return
		}
		bp.AddPoint(pt)
	}
	err = cw.client.Write(bp)
	if err != nil {
		log.Error(ctx, err, "Error writing influx batch")
		return
	}
}

// HandleDataEvent is called on every new DataMessage from the collector and adds
// the new DataMessage to the buffer of DataMessages and flushes the buffer if
// it is full
func (cw *InfluxWriter) HandleDataEvent(ctx context.Context, dataEvent events.DataEvent) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InfluxWriter/HandleDataEvent")
	defer span.Finish()
	cw.deBuffer.AddData(ctx, dataEvent)
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

	cw.deBuffer.Flush(ctx)
	cw.client.Close()
	cw.runName = ""
}
