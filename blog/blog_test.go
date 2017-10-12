package blog

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

var TEST_STRINGS = [][]byte{[]byte("hello"), []byte("i am a packet"), []byte("im another packet")}

type BlogTest struct {
	Packets           [][]byte
	PacketsToRead     int
	ExpectedNextError error
	ExpectedReadError error
}

var BlogTests = []BlogTest{
	BlogTest{
		Packets:           TEST_STRINGS,
		PacketsToRead:     len(TEST_STRINGS),
		ExpectedNextError: nil,
		ExpectedReadError: nil,
	},
	BlogTest{
		Packets:           TEST_STRINGS,
		PacketsToRead:     len(TEST_STRINGS) + 1,
		ExpectedNextError: io.EOF,
		ExpectedReadError: nil,
	},
}

func checkError(t *testing.T, expected error, actual error) {
	if expected != actual {
		t.Errorf("Expected error, %v, but got error %v instead", expected, actual)
	}
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

		for i := 0; i < blogTest.PacketsToRead; i++ {
			err := rdr.Next()
			if err != nil {
				checkError(t, blogTest.ExpectedNextError, err)
				break
			}
			readPacket, err := ioutil.ReadAll(rdr)
			if err != nil {
				checkError(t, blogTest.ExpectedReadError, err)
				break
			}

			if !reflect.DeepEqual(blogTest.Packets[i], readPacket) {
				t.Errorf("Output packet, %v, does not match input packet %s", readPacket, blogTest.Packets[i])
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
