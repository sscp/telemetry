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
	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"
	sscpproto "github.com/sscp/naturallight-telemetry/proto"
	"log"
	"sync"
	"time"
)

// Collector recieves packets from a packetSource and delivers binary packets
// to all BinaryHandlers and deserialized Proto stucts to Datahandlers.
//
// Collector up a channel and a goroutine for each BinaryHandler and
// DataHandler. A goroutine running processPackets handles the delivery of
// packets and deserialized data to all of the handlers.
type Collector struct {
	packetSource     PacketSource
	binaryHandlers   []BinaryHandler
	binaryChans      []chan ContextPacket
	dataHandlers     []DataHandler
	dataChans        []chan ContextDataMessage
	packetsProcessed int
	waitGroup        *sync.WaitGroup
}

// ContextDataMessage is an internal type for passing packet-scopped context
// though a channel along with the DataMessage pointer
type ContextDataMessage struct {
	ctx  context.Context
	data *sscpproto.DataMessage
}

// NewUDPCollector creates a new Collector that listens on the UDP port
// specified and writes .csv and .blog files
func NewUDPCollector(port int) *Collector {
	ps, err := NewUDPPacketSource(port)
	if err != nil {
		panic(err)
	}
	csvHandler := NewCSVWriter()
	blogHandler := NewBlogWriter()
	return NewCollector(ps, []BinaryHandler{blogHandler}, []DataHandler{csvHandler}, 10)
}

// NewCollector creates a new instance of Collector that reads packets from the
// given PacketSource, and outputs data to the given BinaryHandlers and
// Datahandlers. Channels are setup for each handler with the given bufferSize.
func NewCollector(ps PacketSource, bh []BinaryHandler, dh []DataHandler, bufferSize int) *Collector {
	col := &Collector{
		packetSource:     ps,
		binaryHandlers:   bh,
		dataHandlers:     dh,
		packetsProcessed: 0,
		waitGroup:        &sync.WaitGroup{},
	}
	col.createChannels(bufferSize)
	return col
}

// RecordRun starts listening for and processing packets from the the
// PacketSource
func (col *Collector) RecordRun(ctx context.Context, runName string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "collector/RecordRun")
	defer span.Finish()

	currentTime := time.Now()
	col.startHandlers(ctx, runName, currentTime)
	go col.processPackets()
	col.packetSource.Listen()
}

// GetPacketsProcessed returns the count of the number of packets that
// collector has processed. This count is updated every time a packet has been
// delivered to all BinaryHandlers and all DataHandlers.
func (col *Collector) GetPacketsProcessed() int {
	return col.packetsProcessed
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
}

// createChannels creates an array of channels that holds channels for each
// BinaryHandler in binaryHandlers and an array of channels for each
// DataHandler in dataHandlers.
func (col *Collector) createChannels(bufferSize int) {
	col.binaryChans = make([]chan ContextPacket, len(col.binaryHandlers))
	for i := 0; i < len(col.binaryChans); i++ {
		col.binaryChans[i] = make(chan ContextPacket, bufferSize)
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
			col.binaryChans[i] <- ContextPacket{
				ctx:    ctx,
				packet: packet,
			}
		} else {
			col.binaryHandlers[i].HandleDroppedPacket(ctx)
		}
	}

	// Deserialize ProtoBuf
	dMsg := sscpproto.DataMessage{}
	err := proto.Unmarshal(packet, &dMsg)
	if err != nil {
		log.Print(err)
	}

	// Pass off deserialized data to channels
	for i := 0; i < len(col.dataChans); i++ {
		if len(col.dataChans[i]) != cap(col.dataChans[i]) {
			col.dataChans[i] <- ContextDataMessage{
				ctx:  ctx,
				data: &dMsg,
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
		col.processPacket(ctxPacket.ctx, ctxPacket.packet)
		col.packetsProcessed++
	}
	col.closeChannels()
}

// wrapBinaryHandler creates a goroutine that listens to the given channel and
// calls the BinaryHandler on each packet and context. One is added to the
// given WaitGroup and when the goroutine exits, one is subtracted from the
// WaitGroup.
func wrapBinaryHandler(binaryFunc func(context.Context, []byte), packetChan <-chan ContextPacket, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for packet := range packetChan {
			binaryFunc(packet.ctx, packet.packet)
		}
	}()
}

// wrapDataHandler creates a goroutine that listens to the given channel and
// calls the DataHandler on each packet and context. One is added to the
// given WaitGroup and when the goroutine exits, one is subtracted from the
// WaitGroup.
func wrapDataHandler(dataFunc func(context.Context, *sscpproto.DataMessage), dataMsgChan <-chan ContextDataMessage, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for dataMsg := range dataMsgChan {
			dataFunc(dataMsg.ctx, dataMsg.data)
		}
	}()
}
