package handlers

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	internalproto "github.com/sscp/telemetry/collector/internalproto"
	"github.com/sscp/telemetry/collector/sources"

	"github.com/gocarina/gocsv"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func BenchmarkCSVWriter(b *testing.B) {
	cw, err := NewCSVWriter(CSVConfig{Folder: "."})
	if err != nil {
		b.Fatalf("Could not create CSVWriter: %v", err)
	}
	runName := "bench"
	runStart := time.Now()
	ctx := context.TODO()
	cw.HandleStartRun(ctx, runName, runStart)
	defer cw.HandleEndRun(ctx, time.Now())
	defer os.Remove(GetCSVFileName(runName, runStart))
	zdm := sources.CreateZeroDataMessage()
	b.ResetTimer()
	// run b.N times
	for n := 0; n < b.N; n++ {
		cw.HandleData(ctx, zdm)
	}
	b.StopTimer()
}

func runCSVWriterTest(t *testing.T, numPackets int) {
	testPackets := make([]*internalproto.DataMessage, numPackets)
	for i := range testPackets {
		testPackets[i] = sources.CreateZeroDataMessage()
		num := int32(i)
		testPackets[i].PowerSaveOn = &num
	}

	cw, err := NewCSVWriter(CSVConfig{Folder: "."})
	if err != nil {
		t.Fatalf("Could not create CSVWriter: %v", err)
	}
	runName := "bench"
	runStart := time.Now()
	ctx := context.TODO()

	cw.HandleStartRun(ctx, runName, runStart)
	defer os.Remove(GetCSVFileName(runName, runStart))
	for _, packet := range testPackets {
		cw.HandleData(ctx, packet)
	}
	cw.HandleEndRun(ctx, time.Now())

	csv, err := ioutil.ReadFile(GetCSVFileName(runName, runStart))
	if err != nil {
		t.Errorf("Could not read written file: %v", err)
	}
	csvStr := string(csv)

	csvContent, err := gocsv.MarshalString(&testPackets)

	if strings.TrimSpace(csvStr) != strings.TrimSpace(csvContent) {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(strings.TrimSpace(csvStr), strings.TrimSpace(csvContent), false)

		t.Errorf("CSV mismatch:\n%v", dmp.DiffPrettyText(diffs))
	}

}

func TestCSVWriter(t *testing.T) {
	runCSVWriterTest(t, 30)
}
