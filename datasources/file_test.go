package datasources

import (
	"bytes"
	"github.com/sscp/naturallight-telemetry/blog"
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
		Delay:   0,
	},
	BlogReaderSourceTest{
		Packets: [][]byte{[]byte("hello"), []byte("i am a packet"), []byte("im another packet")},
		Delay:   time.Second / 100,
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
		rdr := ReadPackets(bufRead, blogTest.Delay)

		var i int = 0
		for packet := range rdr.Packets() {
			if !reflect.DeepEqual(blogTest.Packets[i], packet) {
				t.Errorf("Output packet, %v, does not match input packet %s", packet, blogTest.Packets[i])
			}
			i++
		}
	}
}
