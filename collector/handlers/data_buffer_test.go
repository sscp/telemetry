package handlers

import (
	"context"
	"reflect"
	"testing"
	"time"

	internalproto "github.com/sscp/telemetry/collector/internalproto"
	"github.com/sscp/telemetry/collector/sources"
)

func TestDataMessageBuffer(t *testing.T) {
	runDataMessageTest(t, 20, 10)
	runDataMessageTest(t, 20, 1)
	runDataMessageTest(t, 20, 2)
	runDataMessageTest(t, 20, 5)
}

func runDataMessageTest(t *testing.T, numItems int, bufferSize int) {

	in := createTestDataMessages(numItems)
	var out []*internalproto.DataMessage

	dmb := NewDataMessageBuffer(func(ctx context.Context, data []*internalproto.DataMessage) {
		for _, dm := range data {
			out = append(out, dm)
		}
	}, bufferSize)

	for _, dm := range in {
		dmb.AddData(context.TODO(), dm)
	}
	dmb.Flush(context.TODO())

	if len(in) != len(out) {
		t.Errorf("Not all data made it though. In: %v Out: %v", len(in), len(out))
	}
	inIndices := getIndexList(in)
	outIndices := getIndexList(out)
	if !reflect.DeepEqual(inIndices, outIndices) {
		t.Errorf("Buffer corrupted data. In: %v Out: %v", inIndices, outIndices)
	}

	inTimes := getTimeList(in)
	outTimes := getTimeList(out)
	if !reflect.DeepEqual(inTimes, outTimes) {
		t.Errorf("Buffer corrupted data. In: %v Out: %v", inTimes, outTimes)
	}

}

func createTestDataMessages(numItems int) []*internalproto.DataMessage {
	var arr []*internalproto.DataMessage
	for i := 0; i < numItems; i++ {
		zdm := sources.CreateZeroDataMessage()
		time := time.Now().UnixNano()
		zdm.TimeCollected = &time
		index := uint32(i)
		zdm.RegenEnabled = &index
		arr = append(arr, zdm)
	}
	return arr
}

func getTimeList(dms []*internalproto.DataMessage) []int64 {
	times := make([]int64, len(dms))
	for i, dm := range dms {
		times[i] = dm.GetTimeCollected()
	}
	return times
}

func getIndexList(dms []*internalproto.DataMessage) []uint32 {
	indies := make([]uint32, len(dms))
	for i, dm := range dms {
		indies[i] = dm.GetRegenEnabled()
	}
	return indies
}
