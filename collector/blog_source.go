package collector

import (
	"context"
	"github.com/sscp/telemetry/blog"
	"io"
	"time"
)

// BlogPacketSource is a PacketSource that reads from an io.Reader
// The delay between packets can be set to some constant
type BlogPacketSource struct {
	reader   io.Reader
	doneChan chan bool
	outChan  chan *ContextPacket
}

// NewBlogPacketSource instantiates a BlogPacketSource
// It reads packets to the specified output channel and waits the given
// duration between reading packets
func NewBlogPacketSource(r io.Reader, d time.Duration) PacketSource {
	return &BlogPacketSource{
		reader:   r,
		doneChan: make(chan bool),
		outChan:  make(chan *ContextPacket),
	}
}

// Listen reads packets from the file sequentially until the file is empty, then calls Close
func (bps *BlogPacketSource) Listen() {
	rdr := blog.NewReader(bps.reader)
	go func() {
		for {
			readPacket, err := rdr.NextPacket()
			if err != nil {
				if err == io.EOF {
					bps.Close()
					bps.doneChan <- true
					break
				} else {
					panic(err)
				}
			}
			bps.outChan <- &ContextPacket{
				ctx:    context.TODO(),
				packet: readPacket,
			}
		}
	}()
}

// Packets returns the channel into which all the read packets are placed
func (bps *BlogPacketSource) Packets() <-chan *ContextPacket {
	return bps.outChan
}

// Close closes the Packets channel
//
// This is called when the end of the stream is reached to wait until the
// goroutine exits and there are no more packets
func (bps *BlogPacketSource) Close() {
	// Wait on the done channel
	<-bps.doneChan
}
