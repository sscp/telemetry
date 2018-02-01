package collector

import (
	"github.com/golang/protobuf/proto"
	sscpproto "github.com/sscp/naturallight-telemetry/proto"

	"context"
	//"github.com/opentracing/opentracing-go"
	"testing"
	"time"
)

type CollectorTest struct {
	PacketsPerSecond   int
	BufferSize         int
	TestTime           time.Duration
	BinaryHandlerDelay time.Duration
	DataHandlerDelay   time.Duration
	BinaryReceiveAll   bool
	DataReceiveAll     bool
}

type testBinaryHandler struct {
	DeliveryCount  int
	DroppedPackets int
	t              *testing.T
	delay          time.Duration
}

func newTestBinaryHandler(t *testing.T, delay time.Duration) *testBinaryHandler {
	return &testBinaryHandler{
		t:     t,
		delay: delay,
	}
}

func (tbh *testBinaryHandler) HandleStartRun(ctx context.Context, runName string, startTime time.Time) {
	if runName == "" {
		tbh.t.Errorf("runName from collector empty")
	}
	if startTime.Sub(time.Now()) > 10*time.Millisecond {
		tbh.t.Errorf("startTime from collector invalid")
	}
}

func (tbh *testBinaryHandler) HandleEndRun(ctx context.Context, endTime time.Time) {
	if endTime.Sub(time.Now()) > 10*time.Millisecond {
		tbh.t.Errorf("endTime from telemetry invalid")
	}
}

func (tbh *testBinaryHandler) HandlePacket(ctx context.Context, packet []byte) {
	dm := sscpproto.DataMessage{}
	err := proto.Unmarshal(packet, &dm)
	if err != nil {
		tbh.t.Errorf("Expected no deserialiation error, instead got %v", err)
	}
	if dm.GetMotorControllerSpeed() != float32(0.0) {
		tbh.t.Errorf("Expected struct to be zeroed but got for motor speed %v instead", dm.GetMotorControllerSpeed())
	}
	time.Sleep(tbh.delay)
	tbh.DeliveryCount++
}

func (tbh *testBinaryHandler) HandleDroppedPacket(ctx context.Context) {
	tbh.DroppedPackets++
}

type testDataHandler struct {
	DeliveryCount  int
	DroppedPackets int
	delay          time.Duration
	t              *testing.T
}

func newTestDataHandler(t *testing.T, delay time.Duration) *testDataHandler {
	return &testDataHandler{
		t:     t,
		delay: delay,
	}
}

func (tdh *testDataHandler) HandleStartRun(ctx context.Context, runName string, startTime time.Time) {
	if runName == "" {
		tdh.t.Errorf("runName from collector empty")
	}
	if startTime.Sub(time.Now()) > 10*time.Millisecond {
		tdh.t.Errorf("startTime from collector invalid")
	}
}

func (tdh *testDataHandler) HandleEndRun(ctx context.Context, endTime time.Time) {
	if endTime.Sub(time.Now()) > 10*time.Millisecond {
		tdh.t.Errorf("endTime from collector invalid")
	}
}

func (tbh *testDataHandler) HandleData(ctx context.Context, data *sscpproto.DataMessage) {
	time.Sleep(tbh.delay)
	tbh.DeliveryCount++
}

func (tbh *testDataHandler) HandleDroppedData(ctx context.Context) {
	tbh.DroppedPackets++
}

func runCollectorTest(t *testing.T, test CollectorTest) {
	bh := newTestBinaryHandler(t, test.BinaryHandlerDelay)
	dh := newTestDataHandler(t, test.DataHandlerDelay)
	zps := NewZeroPacketSource(test.PacketsPerSecond)
	telem := NewCollector(zps, []BinaryHandler{BinaryHandler(bh)}, []DataHandler{DataHandler(dh)}, test.BufferSize)
	ctx := context.TODO()
	telem.RecordRun(ctx, "test")
	time.Sleep(test.TestTime)
	telem.Close(ctx)
	expectedPackets := int(float64(test.PacketsPerSecond) * test.TestTime.Seconds())
	if telem.GetPacketsProcessed() < expectedPackets {
		t.Errorf("Expected to process %v packets, but collector only processed %v packets", expectedPackets, telem.GetPacketsProcessed())
	}
	if test.BinaryReceiveAll {
		if bh.DeliveryCount < telem.GetPacketsProcessed() {
			t.Errorf("Expected all packets to be delivered to binary handler, but %v packets were processed and %v delivered", telem.GetPacketsProcessed(), bh.DeliveryCount)
		}
	} else {
		if bh.DeliveryCount == telem.GetPacketsProcessed() {
			t.Errorf("Expected binary handler to fall behind, but %v packets were processed and %v delivered", telem.GetPacketsProcessed(), bh.DeliveryCount)

		}

	}

	if test.DataReceiveAll {
		if dh.DeliveryCount < telem.GetPacketsProcessed() {
			t.Errorf("Expected all packets to be delivered to data handler, but %v packets were processed and %v delivered", telem.GetPacketsProcessed(), dh.DeliveryCount)
		}
	} else {
		if dh.DeliveryCount == telem.GetPacketsProcessed() {
			t.Errorf("Expected data handler to fall behind, but %v packets were processed and %v delivered", telem.GetPacketsProcessed(), dh.DeliveryCount)

		}

	}
}

var collectorTests []CollectorTest = []CollectorTest{
	CollectorTest{
		PacketsPerSecond:   100,
		BufferSize:         1,
		TestTime:           500 * time.Millisecond,
		BinaryHandlerDelay: 5 * time.Millisecond,
		DataHandlerDelay:   5 * time.Millisecond,
		BinaryReceiveAll:   true,
		DataReceiveAll:     true,
	},
	CollectorTest{
		PacketsPerSecond:   100,
		BufferSize:         1,
		TestTime:           500 * time.Millisecond,
		BinaryHandlerDelay: 100 * time.Millisecond,
		DataHandlerDelay:   5 * time.Millisecond,
		BinaryReceiveAll:   false,
		DataReceiveAll:     true,
	},
	CollectorTest{
		PacketsPerSecond:   100,
		BufferSize:         1,
		TestTime:           500 * time.Millisecond,
		BinaryHandlerDelay: 5 * time.Millisecond,
		DataHandlerDelay:   100 * time.Millisecond,
		BinaryReceiveAll:   true,
		DataReceiveAll:     false,
	},

	CollectorTest{
		PacketsPerSecond:   100,
		BufferSize:         10,
		TestTime:           500 * time.Millisecond,
		BinaryHandlerDelay: 5 * time.Millisecond,
		DataHandlerDelay:   5 * time.Millisecond,
		BinaryReceiveAll:   true,
		DataReceiveAll:     true,
	},
	CollectorTest{
		PacketsPerSecond:   250,
		BufferSize:         10,
		TestTime:           500 * time.Millisecond,
		BinaryHandlerDelay: 1 * time.Millisecond,
		DataHandlerDelay:   1 * time.Millisecond,
		BinaryReceiveAll:   true,
		DataReceiveAll:     true,
	},
	CollectorTest{
		PacketsPerSecond:   250,
		BufferSize:         10,
		TestTime:           500 * time.Millisecond,
		BinaryHandlerDelay: 0 * time.Millisecond,
		DataHandlerDelay:   0 * time.Millisecond,
		BinaryReceiveAll:   true,
		DataReceiveAll:     true,
	},
	CollectorTest{
		PacketsPerSecond:   250,
		BufferSize:         10,
		TestTime:           500 * time.Millisecond,
		BinaryHandlerDelay: 10 * time.Millisecond,
		DataHandlerDelay:   10 * time.Millisecond,
		BinaryReceiveAll:   false,
		DataReceiveAll:     false,
	},
}

func TestCollector(t *testing.T) {
	for _, test := range collectorTests {
		runCollectorTest(t, test)
	}
}
