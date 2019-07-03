package sources

import "github.com/sscp/telemetry/events"

// RawEventSource abstracts over a source of data packets, can be a file or
// listening for UDP packets
type RawEventSource interface {
	// RawEvents returns a reference to the output channel of rawEvents
	// produced by the DataSource
	RawEvents() <-chan *events.ContextRawEvent

	// Listen begins collecting packets and putting them on the output
	// channel.
	Listen()

	// Close stops putting packets on the output channel
	Close() error
}
