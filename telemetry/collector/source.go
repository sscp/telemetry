package collector

import (
	"context"
)

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
