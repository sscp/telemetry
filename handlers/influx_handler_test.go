package handlers

import (
	"context"
	"testing"
	"time"
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

	testEvents := createTestDataEvents(b.N)

	// run b.N times
	for n := 0; n < b.N; n++ {
		iw.HandleDataEvent(ctx, testEvents[n])
	}
	b.StopTimer()
}
