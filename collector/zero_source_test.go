package collector

import (
	"github.com/golang/protobuf/proto"
	sundaeproto "github.com/sscp/telemetry/collector/sundae"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZeroPacketSource(t *testing.T) {
	zps := NewZeroPacketSource(1000)

	zps.Listen()

	for i := 0; i < 10; i++ {
		ctxPacket := <-zps.Packets()
		dm := sundaeproto.DataMessage{}
		err := proto.Unmarshal(ctxPacket.packet, &dm)
		assert.Nil(t, err)
		assert.Equal(
			t, dm.GetMotorControllerSpeed(), float32(0.0),
			"Motor controller speed should be zero",
		)
	}
}
