// Package collector contains all code that directly collects and processes the
// car data stream. There are three primary concepts in addition to the core
// collector defined below:
//
// PacketSource (defined in source.go) which collects packets from UDP or
// some other source and passes them to the core collector. Implementations of a
// UDPPacketSource and a BlogPacketSource can be found in udp_source.go and
// blog_source.go respectively.
//
// DataHandler (defined in handlers.go) is a sink for deserialized data in
// the form of DataMessage. The implementation of a CSVWriter can be found in
// csv_handler.go.
//
// BinaryHandler (defined in handlers.go) is a sink for binary data in the form
// of byte arrays. The implementation of a BlogWriter can be found in
// blog_handler.go
//
// TODO(jbeasley) 2017-01-26 - collector needs a logging system so that there
// are no calls to panic()
package collector

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/sscp/telemetry/collector/handlers"
	internalproto "github.com/sscp/telemetry/collector/internalproto"
	"github.com/sscp/telemetry/collector/sources"
	sundaeproto "github.com/sscp/telemetry/collector/sundae"
	"github.com/sscp/telemetry/log"

	"github.com/opentracing/opentracing-go"
)

//go:generate protoc --go_out=internalproto ./internalproto/data_message.proto
//go:generate protoc-go-inject-tag -input=./internalproto/data_message.pb.go

const defaultBufferSize = 10

// Collector receives packets from a packetSource and delivers binary packets
// to all BinaryHandlers and deserialized Proto stucts to Datahandlers.
//
// Collector up a channel and a goroutine for each BinaryHandler and
// DataHandler. A goroutine running processPackets handles the delivery of
// packets and deserialized data to all of the handlers.
type Collector struct {
	packetSource   sources.PacketSource
	binaryHandlers []handlers.BinaryHandler
	binaryChans    []chan sources.ContextPacket
	dataHandlers   []handlers.DataHandler
	dataChans      []chan ContextDataMessage
	waitGroup      *sync.WaitGroup
	status         CollectorStatus
}

// CollectorConfig holds config values needed to create a collector
type CollectorConfig struct {
	Port int
	CSV  *handlers.CSVConfig
	Blog *handlers.BlogConfig

	Influx *handlers.InfluxConfig
}

// CollectorStatus holds variables that pertain to the current status of
// collector. Variables are reset to zero values at then end of a run
type CollectorStatus struct {
	Collecting bool
	RunName    string
	// PacketsProcessed is the count of the number of packets that
	// collector has processed. This count is updated every time a packet
	// has been delivered to all BinaryHandlers and all DataHandlers.
	PacketsProcessed int64
}

// ContextDataMessage is an internal type for passing packet-scoped context
// though a channel along with the DataMessage pointer
type ContextDataMessage struct {
	ctx  context.Context
	data *internalproto.DataMessage
}

// NewUDPCollector creates a new Collector that listens on the UDP port
// specified and writes .csv and .blog files
func NewUDPCollector(cfg CollectorConfig) (*Collector, error) {
	ps, err := sources.NewUDPPacketSource(cfg.Port)
	if err != nil {
		return nil, err
	}

	var binaryHandlers []handlers.BinaryHandler
	var dataHandlers []handlers.DataHandler

	if cfg.CSV != nil {
		csvHandler, err := handlers.NewCSVWriter(*cfg.CSV)
		if err != nil {
			return nil, err
		}
		dataHandlers = append(dataHandlers, csvHandler)
	}

	if cfg.Blog != nil {
		blogHandler, err := handlers.NewBlogWriter(*cfg.Blog)
		if err != nil {
			return nil, err
		}
		binaryHandlers = append(binaryHandlers, blogHandler)
	}
	if cfg.Influx != nil {
		influxHandler, err := handlers.NewInfluxWriter(*cfg.Influx)
		if err != nil {
			return nil, err
		}
		dataHandlers = append(dataHandlers, influxHandler)
	}

	return NewCollector(ps, binaryHandlers, dataHandlers), nil
}

// NewCollector creates a new instance of Collector that reads packets from the
// given PacketSource, and outputs data to the given BinaryHandlers and
// Datahandlers. Channels are setup for each handler with the given bufferSize.
func NewCollector(ps sources.PacketSource, bh []handlers.BinaryHandler, dh []handlers.DataHandler) *Collector {
	col := &Collector{
		packetSource:   ps,
		binaryHandlers: bh,
		dataHandlers:   dh,
		waitGroup:      &sync.WaitGroup{},
		status:         CollectorStatus{},
	}
	return col
}

// RecordRun starts listening for and processing packets from the
// PacketSource
func (col *Collector) RecordRun(ctx context.Context, runName string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "collector/RecordRun")
	defer span.Finish()

	currentTime := time.Now()
	col.createChannels(defaultBufferSize)
	col.startHandlers(ctx, runName, currentTime)
	go col.processPackets()
	// Reset all status variables to their zero values
	col.status = CollectorStatus{}
	col.status.RunName = runName
	col.status.Collecting = true

	col.packetSource.Listen()
}

// GetStatus returns the status struct for the collector
func (col *Collector) GetStatus() *CollectorStatus {
	return &col.status
}

// Close stops listening for packets and waits until the handlers have finished
// processing everything in their channels.
func (col *Collector) Close(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "collector/Close")
	defer span.Finish()

	col.packetSource.Close()
	currentTime := time.Now() // Get time after we stop recording, not processing
	col.waitGroup.Wait()
	col.stopHandlers(ctx, currentTime)
	col.status.Collecting = false

	// This is a good time to garbage collect
	runtime.GC()

}

// createChannels creates an array of channels that holds channels for each
// BinaryHandler in binaryHandlers and an array of channels for each
// DataHandler in dataHandlers.
func (col *Collector) createChannels(bufferSize int) {
	col.binaryChans = make([]chan sources.ContextPacket, len(col.binaryHandlers))
	for i := 0; i < len(col.binaryChans); i++ {
		col.binaryChans[i] = make(chan sources.ContextPacket, bufferSize)
	}
	col.dataChans = make([]chan ContextDataMessage, len(col.dataHandlers))
	for i := 0; i < len(col.dataChans); i++ {
		col.dataChans[i] = make(chan ContextDataMessage, bufferSize)
	}
}

// closeChannels closes all of the dataChans and binaryChans. Should only be
// called after the channels are empty and when both arrays of channels already
// exist.
func (col *Collector) closeChannels() {
	for i := 0; i < len(col.binaryChans); i++ {
		close(col.binaryChans[i])
	}
	for i := 0; i < len(col.dataChans); i++ {
		close(col.dataChans[i])
	}
}

// stopHandlers loops through and calls HandleEndRun on all handlers. Should
// only be called after all data is processed.
func (col *Collector) stopHandlers(ctx context.Context, endTime time.Time) {
	for _, handler := range col.binaryHandlers {
		handler.HandleEndRun(ctx, endTime)
	}
	for _, handler := range col.dataHandlers {
		handler.HandleEndRun(ctx, endTime)
	}

}

// Starts all of the Handler goroutines that listen to the binaryChans and
// dataChans. To be called after the channels are created by createChannels.
func (col *Collector) startHandlers(ctx context.Context, runName string, startTime time.Time) {
	for i, handler := range col.binaryHandlers {
		wrapBinaryHandler(handler.HandlePacket, col.binaryChans[i], col.waitGroup)
		handler.HandleStartRun(ctx, runName, startTime)
	}
	for i, handler := range col.dataHandlers {
		wrapDataHandler(handler.HandleData, col.dataChans[i], col.waitGroup)
		handler.HandleStartRun(ctx, runName, startTime)
	}
}

// processPacket sends a single packet to all binaryChans in binary form and to
// all dataChans in deserialized form
func (col *Collector) processPacket(ctx context.Context, packet []byte) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "collector/processPacket")
	defer span.Finish()

	// Forward the binary packets first
	for i := range col.binaryChans {
		if len(col.binaryChans[i]) != cap(col.binaryChans[i]) {
			col.binaryChans[i] <- sources.ContextPacket{
				Ctx:    ctx,
				Packet: packet,
			}
		} else {
			col.binaryHandlers[i].HandleDroppedPacket(ctx)
		}
	}

	// Deserialize ProtoBuf
	dMsg, err := sundaeproto.Deserialize(ctx, packet)
	if err != nil {
		log.Error(ctx, err, "Could not deserialize protobuf")
		return
	}

	// Pass off deserialized data to channels
	for i := 0; i < len(col.dataChans); i++ {
		if len(col.dataChans[i]) != cap(col.dataChans[i]) {
			col.dataChans[i] <- ContextDataMessage{
				ctx:  ctx,
				data: dMsg,
			}
		} else {
			col.dataHandlers[i].HandleDroppedData(ctx)
		}
	}
}

// Listens to the incomming packets on the DataSource's channel and processes
// them
func (col *Collector) processPackets() {
	for ctxPacket := range col.packetSource.Packets() {
		col.processPacket(ctxPacket.Ctx, ctxPacket.Packet)
		col.status.PacketsProcessed++
	}
	col.closeChannels()
}

// wrapBinaryHandler creates a goroutine that listens to the given channel and
// calls the BinaryHandler on each packet and context. One is added to the
// given WaitGroup and when the goroutine exits, one is subtracted from the
// WaitGroup.
func wrapBinaryHandler(binaryFunc func(context.Context, []byte), packetChan <-chan sources.ContextPacket, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for packet := range packetChan {
			binaryFunc(packet.Ctx, packet.Packet)
		}
	}()
}

// wrapDataHandler creates a goroutine that listens to the given channel and
// calls the DataHandler on each packet and context. One is added to the
// given WaitGroup and when the goroutine exits, one is subtracted from the
// WaitGroup.
func wrapDataHandler(dataFunc func(context.Context, *internalproto.DataMessage), dataMsgChan <-chan ContextDataMessage, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for dataMsg := range dataMsgChan {
			dataFunc(dataMsg.ctx, dataMsg.data)
		}
	}()
}
