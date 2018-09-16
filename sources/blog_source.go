package sources

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/sscp/telemetry/blog"
	"github.com/sscp/telemetry/events"
)

// BlogPacketSource is a PacketSource that reads from an io.Reader
// The delay between packets can be set to some constant
type BlogPacketSource struct {
	reader   io.Reader
	doneChan chan bool
	outChan  chan *events.ContextRawEvent
}

// NewBlogPacketSource instantiates a BlogPacketSource
// It reads packets to the specified output channel and waits the given
// duration between reading packets
func NewBlogPacketSource(r io.Reader, d time.Duration) PacketSource {
	return &BlogPacketSource{
		reader:   r,
		doneChan: make(chan bool),
		outChan:  make(chan *events.ContextRawEvent),
	}
}

// Listen reads packets from the file sequentially until the file is empty, then calls Close
func (bps *BlogPacketSource) Listen() {
	rdr := blog.NewReader(bps.reader)
	for {
		readPacket, err := rdr.NextPacket()
		if err != nil {
			if err == io.EOF {
				bps.Close()
				bps.doneChan <- true
				break
			} else {
				log.Fatal(err)
			}
		}
		// TODO: NewRawDataEvent sets CollectedTimeNanos to
		// current time, maybe try to pull from blog?
		bps.outChan <- &events.ContextRawEvent{
			Context:  context.Background(),
			RawEvent: events.NewRawEventNow(readPacket),
		}
	}
}

// RawEvents returns the channel into which all the read packets are placed
func (bps *BlogPacketSource) RawEvents() <-chan *events.ContextRawEvent {
	return bps.outChan
}

// Close closes the RawEvents channel
//
// This is called when the end of the stream is reached to wait until the
// goroutine exits and there are no more packets
func (bps *BlogPacketSource) Close() {
	// Wait on the done channel
	<-bps.doneChan
}
