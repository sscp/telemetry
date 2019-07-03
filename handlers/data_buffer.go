package handlers

import (
	"context"

	"github.com/sscp/telemetry/events"
)

// DataMapBuffer holds DataMessages before emptying the buffer by calling a
// callback function to write everything out of the buffer. The buffer is
// flushed whenever AddData is called and the buffer becomes full or whenever
// Flush is called. Ensure to call flush to close out the buffer at the end of
// use.
type DataEventBuffer struct {
	dataBuf   []events.DataEvent
	lastIndex int
	writeData func(context.Context, []events.DataEvent)
}

// NewDataMapBuffer constructs a new DataMapBuffer with the given
// callback function and of the given size.
func NewDataEventBuffer(callback func(context.Context, []events.DataEvent), size int) *DataEventBuffer {
	return &DataEventBuffer{
		dataBuf:   make([]events.DataEvent, 10),
		lastIndex: -1,
		writeData: callback,
	}
}

// AddData adds a single DataMessage to the buffer and if after that message is
// added, the buffer is full, it calls Flush to empty the buffer.
func (dmb *DataEventBuffer) AddData(ctx context.Context, dm events.DataEvent) {
	dmb.dataBuf[dmb.lastIndex+1] = dm
	dmb.lastIndex++
	if dmb.lastIndex >= len(dmb.dataBuf)-1 {
		// All data is fresh, write the whole buffer
		dmb.Flush(ctx)
	}
}

// Flush empties the buffer by calling the callback function and passing it all
// the data in the buffer, then reseting the buffer index.
func (deb *DataEventBuffer) Flush(ctx context.Context) {
	if deb.lastIndex != -1 {
		// Slice is exclusive so +1 to lastIndex
		deb.writeData(ctx, deb.dataBuf[0:deb.lastIndex+1])
	}
	deb.lastIndex = -1

}
