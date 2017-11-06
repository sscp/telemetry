package datasources

import (
	"github.com/sscp/naturallight-telemetry/blog"
	"io"
	"time"
)

// BlogReaderSource is a DataSource that reads from an io.Reader
// The delay between packets can be set to some constant
type BlogReaderSource struct {
	reader   io.Reader
	delay    time.Duration
	doneChan chan bool
	outChan  chan []byte
}

// ReadPackets instantiates and starts a BlogReaderSource
// It reads packets to the specified output channel and waits the given
// duration between reading packets
func ReadPackets(r io.Reader, d time.Duration) DataSource {
	brs := &BlogReaderSource{reader: r, delay: d, doneChan: make(chan bool), outChan: make(chan []byte)}
	brs.read()
	return brs
}

// read reads packets from the file sequentially until the file is empty, then calls Close
func (brs *BlogReaderSource) read() {
	rdr := blog.NewReader(brs.reader)
	go func() {
		for {
			readPacket, err := rdr.NextPacket()
			if err != nil {
				if err == io.EOF {
					brs.Close()
					brs.doneChan <- true
					break
				} else {
					panic(err)
				}
			}
			brs.outChan <- readPacket
			time.Sleep(brs.delay)
		}
	}()
}

// Packets returns the channel into which all the read packets are placed
func (brs *BlogReaderSource) Packets() chan []byte {
	return brs.outChan
}

// Close closes the Packets channel
//
// This is called when the end of the stream is reached to wait until the
// goroutine exits and there are no more packets
func (brs *BlogReaderSource) Close() {
	// Wait on the done channel
	<-brs.doneChan
}
