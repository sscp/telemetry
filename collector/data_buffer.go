package collector

import (
	"context"

	sundaeproto "github.com/sscp/telemetry/collector/sundae"
)

// DataMessageBuffer holds DataMessages before emptying the buffer by calling a
// callback function to write everything out of the buffer. The buffer is
// flushed whenever AddData is called and the buffer becomes full or whenever
// Flush is called. Ensure to call flush to close out the buffer at the end of
// use.
type DataMessageBuffer struct {
	dataBuf   []*sundaeproto.DataMessage
	lastIndex int
	writeData func(context.Context, []*sundaeproto.DataMessage)
}

// NewDataMessageBuffer constructs a new DataMessageBuffer with the given
// callback function and of the given size.
func NewDataMessageBuffer(callback func(context.Context, []*sundaeproto.DataMessage), size int) *DataMessageBuffer {
	return &DataMessageBuffer{
		dataBuf:   make([]*sundaeproto.DataMessage, 10),
		lastIndex: -1,
		writeData: callback,
	}
}

// AddData adds a single DataMessage to the buffer and if after that message is
// added, the buffer is full, it calls Flush to empty the buffer.
func (dmb *DataMessageBuffer) AddData(ctx context.Context, dm *sundaeproto.DataMessage) {
	dmb.dataBuf[dmb.lastIndex+1] = dm
	dmb.lastIndex++
	if dmb.lastIndex >= len(dmb.dataBuf)-1 {
		// All data is fresh, write the whole buffer
		dmb.Flush(ctx)
	}
}

// Flush empties the buffer by calling the callback function and passing it all
// the data in the buffer, then reseting the buffer index.
func (dmb *DataMessageBuffer) Flush(ctx context.Context) {
	if dmb.lastIndex != -1 {
		// Slice is exclusive so +1 to lastIndex
		dmb.writeData(ctx, dmb.dataBuf[0:dmb.lastIndex+1])
	}
	dmb.lastIndex = -1

}
