package collector

import (
	sundaeproto "github.com/sscp/telemetry/collector/sundae"
)

type DataMessageBuffer struct {
	dataBuf   []*sundaeproto.DataMessage
	lastIndex int
	writeData func([]*sundaeproto.DataMessage)
}

func NewDataMessageBuffer(callback func([]*sundaeproto.DataMessage), size int) *DataMessageBuffer {
	return &DataMessageBuffer{
		dataBuf:   make([]*sundaeproto.DataMessage, 10),
		lastIndex: -1,
		writeData: callback,
	}
}

func (dmb *DataMessageBuffer) AddData(dm *sundaeproto.DataMessage) {
	dmb.dataBuf[dmb.lastIndex+1] = dm
	dmb.lastIndex++
	if dmb.lastIndex >= len(dmb.dataBuf)-1 {
		// All data is fresh, write the whole buffer
		dmb.Flush()
	}
}

func (dmb *DataMessageBuffer) Flush() {
	if dmb.lastIndex != -1 {
		// Slice is exclusive so +1 to lastIndex
		dmb.writeData(dmb.dataBuf[0 : dmb.lastIndex+1])
	}
	dmb.lastIndex = -1

}
