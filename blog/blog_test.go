package blog

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

type BlogTest struct {
	Packets [][]byte
}

var BlogTests = []BlogTest{
	BlogTest{
		Packets: [][]byte{[]byte("hello"), []byte("i am a packet"), []byte("im another packet")},
	},
}

func TestBlog(t *testing.T) {

	for _, blogTest := range BlogTests {
		buf := new(bytes.Buffer)

		writer := NewWriter(buf)
		for _, s := range blogTest.Packets {
			writer.Write(s)
		}

		bufRead := bytes.NewReader(buf.Bytes())
		rdr := NewReader(bufRead)

		for i := 0; i < len(blogTest.Packets)+1; i++ {
			err := rdr.Next()
			if err != nil {
				if i == len(blogTest.Packets) && err == io.EOF {
					// End of stream
					break
				} else {
					t.Errorf("Unexpected error while advancing packet: %v", err)
				}
			}
			readPacket, err := ioutil.ReadAll(rdr)
			if err != nil {
				t.Errorf("Unexpected error while reading packet %v", err)
			}
			if !reflect.DeepEqual(blogTest.Packets[i], readPacket) {
				t.Errorf("Output packet, %v, does not match input packet %s", readPacket, blogTest.Packets[i])
			}
		}

		// Reset reader
		bufRead.Seek(0, io.SeekStart)

		for i := 0; i < len(blogTest.Packets)+1; i++ {
			packet, err := rdr.NextPacket()
			if err != nil {
				if i == len(blogTest.Packets) && err == io.EOF {
					// End of stream
					break
				} else {
					t.Errorf("Unexpected error while advancing packet: %v", err)
				}

			}

			if !reflect.DeepEqual(blogTest.Packets[i], packet) {
				t.Errorf("Output packet, %v, does not match input packet %s", packet, blogTest.Packets[i])
			}

		}
	}
}

func TestBlogReadZero(t *testing.T) {
	buf := new(bytes.Buffer)

	writer := NewWriter(buf)
	for _, s := range BlogTests[0].Packets {
		writer.Write(s)
	}

	bufRead := bytes.NewReader(buf.Bytes())
	rdr := NewReader(bufRead)
	_, err := ioutil.ReadAll(rdr)
	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}
}
