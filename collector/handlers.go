package collector

import (
	"context"
	"fmt"
	sundaeproto "github.com/sscp/telemetry/collector/sundae"
	"time"
)

// DataHandler is a sink for DataMessages from the collector.
//
// A DataHandler can be plugged into the collector and will recieve all the
// collected DataMessages from the car. Care has been taken to ensure that a
// slow DataHandler design does not impact the overall performance of the
// collector, so if the DataHandler processes packets at a slower rate than the
// collector recieves them, packets will be dropped and the DataHandler
// notified by a call to HandleDroppedData.
type DataHandler interface {
	// HandleStartRun is called when collector starts collecting data and
	// includes the name of the run and the start time of the run. Any
	// run-specific setup should occur in this method, such as creating a
	// file and setting up buffers.
	HandleStartRun(context.Context, string, time.Time)

	// HandleData is called by collector on every incomming packet. This
	// method is performance critical, so if it is slow, the DataHandler
	// will not recieve every packet from collector. This method should be
	// benchmarked to verify that it is fast enough to recieve all data.
	HandleData(context.Context, *sundaeproto.DataMessage)

	// HandleDroppedData is called by collector whenever the DataHandler
	// falls behind and misses a packet. This is a performance critical
	// method and should only be used to collect statistics on the number
	// of packets dropped.
	HandleDroppedData(context.Context)

	// HandleEndRun is called by collector when data collection stops and
	// all data is processed. It should be used to flush buffers and close
	// files/connections.
	HandleEndRun(context.Context, time.Time)
}

// BinaryHandler can be plugged into the collector and will recieve all the
// collected packets as a []byte from the car. Care has been taken to ensure
// that a slow BinaryHandler design does not impact the overall performance of
// the collector, so, as with DataHandlers, if the BinaryHandler processes
// packets at a slower rate than the collector recieves them, packets will be
// dropped and the BinaryHandler notified by a call to HandleDroppedPacket.
type BinaryHandler interface {
	// HandleStartRun is called when collector starts collecting data and
	// includes the name of the run and the start time of the run. Any
	// run-specific setup should occur in this method, such as creating a
	// file and setting up buffers.
	HandleStartRun(context.Context, string, time.Time)

	// HandlePacket is called by collector on every incomming packet. This
	// method is performance critical, so if it is slow, the BinaryHandler
	// will not recieve every packet from collector. This method should be
	// benchmarked to verify that it is fast enough to recieve all data.
	HandlePacket(context.Context, []byte)

	// HandleDroppedPacket is called by collector whenever the
	// BinaryHandler falls behind and misses a packet. This is a
	// performance critical method and should only be used to collect
	// statistics on the number of packets dropped.
	HandleDroppedPacket(context.Context)

	// HandleEndRun is called by collector when data collection stops and
	// all data is processed. It should be used to flush buffers and close
	// files/connections.
	HandleEndRun(context.Context, time.Time)
}

// GetBaseFileName returns a unique string for each run that combines the run
// name with the current time.
func GetBaseFileName(runName string, startTime time.Time) string {
	date := startTime.Format("2006-01-02-15:04:05")
	return fmt.Sprintf("%v_%v", runName, date)
}
