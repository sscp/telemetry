package binary_log

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

var TEST_STRINGS = [][]byte{[]byte("hello"), []byte("i am a packet"), []byte("im another packet")}

func TestBlog(t *testing.T) {
	buf := new(bytes.Buffer)

	enc := NewEncoder(buf)
	for _, s := range TEST_STRINGS {
		enc.Encode(s)
	}

	bufRead := bytes.NewReader(buf.Bytes())
	rdr := NewReader(bufRead)

	for _, s := range TEST_STRINGS {
		if err := rdr.Next(); err != nil {
			panic(err)
		}
		readPacket, err := ioutil.ReadAll(rdr)
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
		fmt.Println(string(readPacket))
		if !reflect.DeepEqual(s, readPacket) {
			t.Errorf("Output packet, %v, does not match input packet %s", readPacket, s)
		}
	}

}
