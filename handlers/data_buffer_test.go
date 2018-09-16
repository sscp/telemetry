package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sscp/telemetry/events"
)

func TestDataEventBuffer(t *testing.T) {
	runDataEventBufferTest(t, 20, 10)
	runDataEventBufferTest(t, 20, 1)
	runDataEventBufferTest(t, 20, 2)
	runDataEventBufferTest(t, 20, 5)
}

func runDataEventBufferTest(t *testing.T, numItems int, bufferSize int) {

	in := createTestDataEvents(numItems)
	var out []events.DataEvent

	dmb := NewDataEventBuffer(func(ctx context.Context, data []events.DataEvent) {
		for _, dm := range data {
			out = append(out, dm)
		}
	}, bufferSize)

	for _, dm := range in {
		dmb.AddData(context.TODO(), dm)
	}
	dmb.Flush(context.TODO())

	assert.Equal(t, in, out, "Buffer should output the exact same data input")
}

func createTestDataEvents(numItems int) []events.DataEvent {
	testEvents := make([]events.DataEvent, numItems)
	for i := 0; i < numItems; i++ {
		testEvents[i] = createTestDataEvent()
	}
	return testEvents
}

func createTestDataEvent() events.DataEvent {
	return events.DataEvent{
		EventMeta: events.EventMeta{
			CollectedTimeNanos: time.Now().UnixNano(),
		},
		Data: map[string]interface{}{
			"test_uint32_data_point": uint32(32),
		},
	}

}
