package datasources

// DataSource abstracts over a source of packets, can be a file or listening for
// UDP packets
//
// Packets is a channel where raw packets are returned
// Close closes the channel, but the channel may close by itself if it reaches
// the end of the file, or there is a natural end to the stream
type DataSource interface {
	Packets() chan []byte
	Close()
}
