package handlers

import (
	"context"
)

// DataMapBuffer holds DataMessages before emptying the buffer by calling a
// callback function to write everything out of the buffer. The buffer is
// flushed whenever AddData is called and the buffer becomes full or whenever
// Flush is called. Ensure to call flush to close out the buffer at the end of
// use.
type DataMapBuffer struct {
	dataBuf   []map[string]interface{}
	lastIndex int
	writeData func(context.Context, []map[string]interface{})
}

// NewDataMapBuffer constructs a new DataMapBuffer with the given
// callback function and of the given size.
func NewDataMapBuffer(callback func(context.Context, []map[string]interface{}), size int) *DataMapBuffer {
	return &DataMapBuffer{
		dataBuf:   make([]map[string]interface{}, 10),
		lastIndex: -1,
		writeData: callback,
	}
}

// AddData adds a single DataMessage to the buffer and if after that message is
// added, the buffer is full, it calls Flush to empty the buffer.
func (dmb *DataMapBuffer) AddData(ctx context.Context, dm map[string]interface{}) {
	dmb.dataBuf[dmb.lastIndex+1] = dm
	dmb.lastIndex++
	if dmb.lastIndex >= len(dmb.dataBuf)-1 {
		// All data is fresh, write the whole buffer
		dmb.Flush(ctx)
	}
}

// Flush empties the buffer by calling the callback function and passing it all
// the data in the buffer, then reseting the buffer index.
func (dmb *DataMapBuffer) Flush(ctx context.Context) {
	if dmb.lastIndex != -1 {
		// Slice is exclusive so +1 to lastIndex
		dmb.writeData(ctx, dmb.dataBuf[0:dmb.lastIndex+1])
	}
	dmb.lastIndex = -1

}
