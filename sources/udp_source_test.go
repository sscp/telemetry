package sources

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sscp/telemetry/events"
)

func TestUDPSendRecv(t *testing.T) {
	src, err := NewUDPPacketSource(3000)
	if err != nil {
		t.Errorf("Error creating packet source: %v", err)
	}

	eventChan := make(chan *events.ContextRawEvent)
	// Send all the packets in the channel
	go SendEventsAsUDP(eventChan, 3000)

	// Listen for those same packets
	go src.Listen()
	defer src.Close()

	for i := 0; i < 1000; i++ {
		// Make a random packet
		packet := make([]byte, i)
		rand.Read(packet)

		// Packet sent to the send channel
		eventChan <- &events.ContextRawEvent{
			Context:  context.Background(),
			RawEvent: events.NewRawEventNow(packet),
		}

		// Listen for the packet on the recv channel
		rawEvent := <-src.RawEvents()

		// Check that everything made it
		assert.Equal(t, packet, rawEvent.Data)
	}
}
