package handlers

import (
	"context"
	"crypto/rand"
	"os"
	"testing"
	"time"

	"github.com/sscp/telemetry/events"
)

func BenchmarkBlogWriter(b *testing.B) {
	bw, err := NewBlogWriter(BlogConfig{Folder: "."})
	if err != nil {
		b.Fatalf("Could not create BlogWriter: %v", err)
	}
	runName := "bench"
	runStart := time.Now()
	ctx := context.TODO()
	bw.HandleStartRun(ctx, runName, runStart)
	defer bw.HandleEndRun(ctx, time.Now())
	defer os.Remove(GetBlogFileName(runName, runStart))

	testPackets := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		// Make a random packet
		testPackets[i] = make([]byte, 1500)
		rand.Read(testPackets[i])
	}
	b.ResetTimer()
	// run b.N times
	for n := 0; n < b.N; n++ {
		bw.HandleRawEvent(ctx, events.NewRawEventNow(testPackets[n]))
	}
	b.StopTimer()
}
