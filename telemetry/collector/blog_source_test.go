package collector

import (
	"bytes"
	"github.com/sscp/naturallight-telemetry/telemetry/blog"
	"reflect"
	"testing"
	"time"
)

type BlogReaderSourceTest struct {
	Packets [][]byte
	Delay   time.Duration
}

var BlogTests = []BlogReaderSourceTest{
	BlogReaderSourceTest{
		Packets: [][]byte{[]byte("hello"), []byte("i am a packet"), []byte("im another packet")},
	},
}

func testBlogReaderSource(t *testing.T) {
	for _, blogTest := range BlogTests {
		buf := new(bytes.Buffer)

		writer := blog.NewWriter(buf)
		for _, s := range blogTest.Packets {
			writer.Write(s)
		}

		bufRead := bytes.NewReader(buf.Bytes())
		bps := NewBlogPacketSource(bufRead, blogTest.Delay)

		var i int = 0
		for packet := range bps.Packets() {
			if !reflect.DeepEqual(blogTest.Packets[i], packet) {
				t.Errorf("Output packet, %v, does not match input packet %s", packet, blogTest.Packets[i])
			}
			i++
		}
	}
}
