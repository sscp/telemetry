package sources

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/sscp/telemetry/events"
)

// UDPListenTimeout is the time to wait for the next packet
const udpListenTimeout = 100 * time.Millisecond

// UDPPacketSource is a PacketSource that reads from a UDP socket
type UDPPacketSource struct {
	port         int
	outChan      chan *events.ContextRawEvent
	doneChan     chan bool
	conn         *net.UDPConn
	packetBuffer []byte
}

// NewUDPPacketSource constructs a UDPPacketSource that listens on the given
// port for packets
func NewUDPPacketSource(port int) (PacketSource, error) {
	ups := &UDPPacketSource{
		port:         port,
		packetBuffer: make([]byte, 2000), // Max packet size is ~1000
	}
	err := ups.setupForListen()
	if err != nil {
		return nil, err
	}
	return ups, nil
}

// setupForListen creates the UDP connection, begins listening, and creates the
// outChan and doneChan to send out received packets and notifies the goroutine
// to stop listening when done
func (ups *UDPPacketSource) setupForListen() error {
	// Listen to the zero port for IPv4 to catch any packet to that port
	// This will catch broadcast packets from the car
	var err error
	ups.conn, err = net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: ups.port,
	})
	if err != nil {
		return err
	}
	ups.outChan = make(chan *events.ContextRawEvent)
	ups.doneChan = make(chan bool)
	return nil
}

// RawEvents is the stream of packets received from UDP
// It is simply a reference to outChan
func (ups *UDPPacketSource) RawEvents() <-chan *events.ContextRawEvent {
	return ups.outChan
}

// Listen spins up a goroutine that listens for packets until it receives a
// signal on the doneChan, in which case it closes the connection and returns
func (ups *UDPPacketSource) Listen() {
	for {
		select {
		case <-ups.doneChan:
			// Close Conn and shutdown goroutine
			ups.conn.Close()
			return
		default:
			ups.readAndForwardPacket()
		}
	}
}

func (ups *UDPPacketSource) readAndForwardPacket() {
	packet, err := ups.readPacket()
	if netError, ok := err.(net.Error); ok {
		// If timeout error, keep looping
		if !netError.Timeout() {
			// Panic if not a timeout error
			log.Fatal(err)
		}
	} else {
		ups.outChan <- &events.ContextRawEvent{
			Context:  context.Background(),
			RawEvent: events.NewRawEventNow(packet),
		}
	}

}

// readPacket reads a single packet into the packetBuffer, then copies the exact
// packet into a new byte array and returns it.
func (ups *UDPPacketSource) readPacket() ([]byte, error) {
	ups.conn.SetDeadline(time.Now().Add(udpListenTimeout))
	numBytes, _, err := ups.conn.ReadFromUDP(ups.packetBuffer)
	if err != nil {
		return nil, err
	}
	// Make a slice for the exact length of the packet and copy the packet
	// into it
	packet := make([]byte, numBytes)
	copy(packet, ups.packetBuffer)
	return packet, nil
}

// Close sends a done signal on doneChan, closes both doneChan, outChan, then
// resets the UDPPacketSource so that it is ready to be reused
func (ups *UDPPacketSource) Close() {
	ups.doneChan <- true
	close(ups.doneChan)
	close(ups.outChan)
	// Reset
	ups.setupForListen()
}

// SendEventsAsUDP sends all the events from the dataSource to the broadcast
// ip on the given port. Packets are spaced by the given delay duration.
func SendEventsAsUDP(eventChan <-chan *events.ContextRawEvent, port int) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for event := range eventChan {
		conn.Write(event.Data)
	}
}
