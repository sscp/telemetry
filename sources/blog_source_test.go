package sources

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sscp/telemetry/blog"
)

func TestBlogReaderSource(t *testing.T) {
	specs := []struct {
		Packets [][]byte
		Delay   time.Duration
	}{
		{
			Packets: [][]byte{[]byte("hello"), []byte("i am a packet"), []byte("im another packet")},
		},
	}

	for _, tt := range specs {
		buf := new(bytes.Buffer)

		writer := blog.NewWriter(buf)
		for _, s := range tt.Packets {
			_, err := writer.Write(s)
			assert.NoError(t, err)
		}

		bufRead := bytes.NewReader(buf.Bytes())
		bps := NewBlogRawEventSource(bufRead, tt.Delay)
		go bps.Listen()

		var i int = 0
		for packet := range bps.RawEvents() {
			assert.Equal(t, tt.Packets[i], packet)
			i++
		}
	}
}
