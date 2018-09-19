package sources

import (
	"context"
	"fmt"

	"github.com/sscp/telemetry/cars/sundae"
	"github.com/sscp/telemetry/events"

	"golang.org/x/time/rate"
)

// ZeroRawEventSource is a RawEventSource that returns only zeroed out DataMessages
// at a given rate
type ZeroRawEventSource struct {
	outChan  chan *events.ContextRawEvent
	doneChan chan bool
	limiter  *rate.Limiter
}

// RawEvents is the stream of zeroed binary packets
// It is simply a reference to outChan
func (zps *ZeroRawEventSource) RawEvents() <-chan *events.ContextRawEvent {
	return zps.outChan
}

// Listen begins sending zeroed packets to the RawEvents channel.
// It launches a gorountine that sen
func (zps *ZeroRawEventSource) Listen() {
	for {
		err := zps.limiter.Wait(context.TODO())
		if err != nil {
			fmt.Println("too fast")
			continue
		}
		zPacket, _ := sundae.CreateZeroPacket()

		select {
		case <-zps.doneChan:
			return
		default:
			zps.outChan <- &events.ContextRawEvent{
				Context:  context.Background(),
				RawEvent: events.NewRawEventNow(zPacket),
			}
		}
	}
}

// Close sends a close signal on doneChan and closes both doneChan and outChan.
// NOTE: this currently does not reset the ZeroRawEventSource to listen again
func (zps *ZeroRawEventSource) Close() error {
	zps.doneChan <- true
	close(zps.outChan)
	close(zps.doneChan)
	return nil
}

// NewZeroRawEventSource constructs a new ZeroRawEventSource that emits zeroed out
// packets at packetsPerSecond
func NewZeroRawEventSource(packetsPerSecond int) RawEventSource {
	return &ZeroRawEventSource{
		outChan:  make(chan *events.ContextRawEvent),
		doneChan: make(chan bool, 1),
		// Only allow one packet out at a time
		limiter: rate.NewLimiter(rate.Limit(packetsPerSecond), 1),
	}
}
