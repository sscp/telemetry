package blog

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

// A Writer provides sequential writing in the .blog format used for Sundae.
// .blog files are streams of packets of arbitrary and variable length like
// this: [uint16][--------packet-------][uint16][--packet--]
// Call Write to append a packet to the end of the stream.
type Writer struct {
	writer io.Writer
}

// NewWriter creates a new Writer writing to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// Write writes b to the .blog packet stream
// Acts like an io.Writer, but writes packets discretely
func (tw *Writer) Write(b []byte) (int, error) {
	var packetSize uint16 = uint16(binary.Size(b))
	err := binary.Write(tw.writer, binary.BigEndian, packetSize)
	if err != nil {
		return 0, err
	}
	err = binary.Write(tw.writer, binary.BigEndian, b)
	if err != nil {
		return 0, err
	}
	return int(packetSize), nil
}

// A Reader provides sequential access to the binary stream file format used for
// Sundae.
// This format contains a stream of binary packets of arbitrary length.
// The Next method moves to the next packet in the stream (including the first)
// and then it can be treated as an io.Reader to access the packet data.
type Reader struct {
	reader           io.Reader
	nextPacketReader io.Reader
}

// NewReader creates a new Reader reading from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: r}
}

// NextPacket returns the next packet in the file
//
// io.EOF is returned as an error if there are no more packets
func (br *Reader) NextPacket() ([]byte, error) {
	err := br.Next()
	if err != nil {
		return nil, err
	}
	readPacket, err := ioutil.ReadAll(br)
	if err != nil {
		return nil, err
	}
	return readPacket, nil
}

// Next advances the reader to the next packet in the stream.
//
// io.EOF is returned when there are no more packets.
func (br *Reader) Next() error {
	var packetSize uint16
	err := binary.Read(br.reader, binary.BigEndian, &packetSize)
	if err != nil {
		return err // Commonly is EOF at the end of
	}
	br.nextPacketReader = io.LimitReader(br.reader, int64(packetSize))
	return nil
}

// Read reads from the current packet in the file.
// It returns 0, io.EOF when it reaches the end of the packet, or if there are
// no more packets left to read. Calls Next to find the first packet if there is
// no currently selected packet.
func (br *Reader) Read(p []byte) (int, error) {
	if br.nextPacketReader == nil {
		err := br.Next()
		if err != nil {
			return 0, err
		}
	}
	return br.nextPacketReader.Read(p)
}
