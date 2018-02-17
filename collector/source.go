package collector

import (
	"context"
	"time"
)

type contextKey string

func (c contextKey) String() string {
	return "packetsource" + string(c)
}

var (
	contextKeyRecievedTime = contextKey("recievedTime")
)

// RecievedTimeFromContext returns the recievedTime recorded by packetSource as well as a bool that is true only if there is a time in the context
func RecievedTimeFromContext(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(contextKeyRecievedTime).(time.Time)
	return t, ok
}

// ContextPacket holds context from
type ContextPacket struct {
	ctx    context.Context
	packet []byte
}

// PacketSource abstracts over a source of data packets, can be a file or
// listening for UDP packets
//
// Packets is a channel where raw packets are returned
// Close closes the channel, but the channel may close by itself if it reaches
// the end of the file, or there is a natural end to the stream
type PacketSource interface {
	// Packets returns a reference to the output channel of packets
	// produced by the DataSource
	Packets() <-chan *ContextPacket

	// Listen begins collecting packets and putting them on the output
	// channel.
	Listen()

	// Close stops putting packets on the output channel
	Close()
}
