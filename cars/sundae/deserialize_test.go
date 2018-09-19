package sundae

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/sscp/telemetry/blog"
	"github.com/sscp/telemetry/events"
)

func runTestOnDataFile(t *testing.T, filename string) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	blogfile, err := os.Open(filename)
	if err != nil {
		t.Errorf("Error opening file (%v): %v", unpaddedTestFilename, err)
	}
	rdr := blog.NewReader(blogfile)
	blogerr := rdr.Next()
	if blogerr != nil {
		t.Errorf("Error reading blog file: %v", err)
	}
	for blogerr == nil {
		packet, err := ioutil.ReadAll(rdr)
		if err != nil {
			t.Errorf("Error reading blog file: %v", err)
		}
		ctx := context.Background()
		_, err = Deserialize(ctx, events.NewRawEventNow(packet))
		if err != nil {
			t.Errorf("Could not deserialize packet: %v", err)
		}
		blogerr = rdr.Next()
	}
}

const unpaddedTestFilename string = "unpadded_test_data.blog"

func TestUnpaddedDeserialize(t *testing.T) {
	runTestOnDataFile(t, unpaddedTestFilename)
}

// This does not work yet... seems to be some other issue in those blog files...
//const paddedTestFilename string = "padded_test_data.blog"
//func TestPaddedDeserialize(t *testing.T) {
//	runTestOnDataFile(t, paddedTestFilename)
//}
