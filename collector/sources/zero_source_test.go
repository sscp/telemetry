package sources

import (
	"github.com/golang/protobuf/proto"
	internalproto "github.com/sscp/telemetry/collector/internalproto"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZeroPacketSource(t *testing.T) {
	zps := NewZeroPacketSource(1000)

	zps.Listen()

	for i := 0; i < 10; i++ {
		ctxPacket := <-zps.Packets()
		dm := internalproto.DataMessage{}
		err := proto.Unmarshal(ctxPacket.Packet, &dm)
		assert.Nil(t, err)
		assert.Equal(
			t, dm.GetMotorControllerSpeed(), float32(0.0),
			"Motor controller speed should be zero",
		)
	}
}
