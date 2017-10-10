package binary_log

import (
	"encoding/binary"
	"io"
)

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{writer: w}
}

func (enc *Encoder) Encode(p interface{}) {
	var packetSize uint16 = uint16(binary.Size(p))
	binary.Write(enc.writer, binary.BigEndian, packetSize)
	binary.Write(enc.writer, binary.BigEndian, p)
}

type Reader struct {
	reader           io.Reader
	nextPacketReader io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{reader: r}
}

func (br *Reader) Next() error {
	var packetSize uint16
	err := binary.Read(br.reader, binary.BigEndian, &packetSize)
	if err != nil {
		// TODO (jbeasley): Differentiate between different error types
		return err
	}
	br.nextPacketReader = io.LimitReader(br.reader, int64(packetSize))
	return nil
}

func (br *Reader) Read(p []byte) (int, error) {
	return br.nextPacketReader.Read(p)
}
