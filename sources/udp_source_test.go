package sources

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"

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
	packets := make([][]byte, 100)

	for i := 0; i < len(packets); i++ {
		// Make a random packet
		packets[i] = make([]byte, i)
		rand.Read(packets[i])
	}
	limiter := rate.NewLimiter(rate.Limit(100), 1)
	for _, packet := range packets {
		err := limiter.Wait(context.TODO())
		if err != nil {
			fmt.Println("f")
		}
		// Packet sent to the send channel
		eventChan <- &events.ContextRawEvent{
			Context:  context.Background(),
			RawEvent: events.NewRawEventNow(packet),
		}
	}
	close(eventChan)

	for i := 0; i < len(packets); i++ {
		// Listen for the packet on the recv channel
		rawEvent := <-src.RawEvents()

		// Check that everything made it
		assert.Equal(t, packets[i], rawEvent.Data)
	}
}
