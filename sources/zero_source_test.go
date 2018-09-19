package sources

import (
	"testing"

	"github.com/sscp/telemetry/cars/sundae"

	"github.com/stretchr/testify/assert"
)

func TestZeroRawEventSource(t *testing.T) {
	zps := NewZeroRawEventSource(1000)
	defer zps.Close()
	go zps.Listen()

	for i := 0; i < 10; i++ {
		rawEvent := <-zps.RawEvents()
		dataEvent, err := sundae.Deserialize(rawEvent.Context, rawEvent.RawEvent)
		assert.Nil(t, err)
		assert.Equal(
			t, dataEvent.Data["motor_controller_speed"], float32(0.0),
			"Motor controller speed should be zero",
		)
	}
}
