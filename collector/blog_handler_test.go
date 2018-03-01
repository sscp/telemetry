package collector

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
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
	defer os.Remove(GetCSVFileName(runName, runStart))
	zdm := CreateZeroDataMessage()
	packet, _ := proto.Marshal(zdm)
	b.ResetTimer()
	// run b.N times
	for n := 0; n < b.N; n++ {
		bw.HandlePacket(ctx, packet)
	}
	b.StopTimer()
}
