package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/sscp/telemetry/collector/sources"
)

func BenchmarkInfluxWriter(b *testing.B) {
	iw, err := NewInfluxWriter(InfluxConfig{
		Addr:     "http://localhost:8086",
		Username: "collector",
		Password: "solarpower",
	})
	if err != nil {
		b.Fatalf("Could not create InfluxWriter: %v", err)
	}
	runName := "bench"
	runStart := time.Now()
	ctx := context.TODO()
	iw.HandleStartRun(ctx, runName, runStart)
	defer iw.HandleEndRun(ctx, time.Now())
	b.ResetTimer()
	// run b.N times
	for n := 0; n < b.N; n++ {
		zdm := sources.CreateRandomDataMessage()
		time := time.Now().UnixNano()
		zdm.TimeCollected = &time
		iw.HandleData(ctx, zdm)
	}
	b.StopTimer()
}
