package sources

import (
	"context"
	"fmt"

	"github.com/sscp/telemetry/cars/sundae"
	"github.com/sscp/telemetry/events"

	"golang.org/x/time/rate"
)

// ZeroPacketSource is a PacketSource that returns only zeroed out DataMessages
// at a given rate
type ZeroPacketSource struct {
	outChan  chan *events.ContextRawEvent
	doneChan chan bool
	limiter  *rate.Limiter
}

// RawEvents is the stream of zeroed binary packets
// It is simply a reference to outChan
func (zps *ZeroPacketSource) RawEvents() <-chan *events.ContextRawEvent {
	return zps.outChan
}

// Listen begins sending zeroed packets to the RawEvents channel.
// It launches a gorountine that sen
func (zps *ZeroPacketSource) Listen() {
	for {
		select {
		case <-zps.doneChan:
			fmt.Println("done")
			return
		default:
			err := zps.limiter.Wait(context.TODO())
			if err != nil {
				fmt.Println("too fast")
				continue
			}
			zPacket, _ := sundae.CreateZeroPacket()

			zps.outChan <- &events.ContextRawEvent{
				Context:  context.Background(),
				RawEvent: events.NewRawEventNow(zPacket),
			}
		}
	}
}

// Close sends a close signal on doneChan and closes both doneChan and outChan.
// NOTE: this currently does not reset the ZeroPacketSource to listen again
func (zps *ZeroPacketSource) Close() {
	zps.doneChan <- true
	<-zps.outChan
	close(zps.doneChan)
	close(zps.outChan)
}

// NewZeroPacketSource constructs a new ZeroPacketSource that emits zeroed out
// packets at packetsPerSecond
func NewZeroPacketSource(packetsPerSecond int) PacketSource {
	return &ZeroPacketSource{
		outChan:  make(chan *events.ContextRawEvent),
		doneChan: make(chan bool, 1),
		// Only allow one packet out at a time
		limiter: rate.NewLimiter(rate.Limit(packetsPerSecond), 1),
	}
}
