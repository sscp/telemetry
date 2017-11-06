package datasources

import (
	"net"
	"time"
)

// Time to wait for the next packet
const UDP_LISTEN_TIMEOUT = 100 * time.Millisecond

// UDPSource is a DataSource that listens for UDP packets sent to the current
// machine's IP
type UDPSource struct {
	outChan      chan []byte
	shutdownChan chan bool
	packetBuffer []byte
	conn         *net.UDPConn
}

// NewUDPSource creates a UDPSource that is bound to the given port.
// It starts listening for packets until Close is called and places the
// packets into an output channel accessible via the Packets function.
func NewUDPSource(port int) UDPSource {
	udpSrc := UDPSource{
		outChan:      make(chan []byte),
		shutdownChan: make(chan bool),
		packetBuffer: make([]byte, 1000), // Max packet size is 512
	}
	udpSrc.listen(port)
	return udpSrc
}

// readPacket reads a single packet from UDP and returns it
// Internally it first reads the packet from UDP into the buffer, then copies
// the packet into a byte array of the exact length of the packet and returns
// it.
func (us *UDPSource) readPacket() ([]byte, error) {
	us.conn.SetDeadline(time.Now().Add(UDP_LISTEN_TIMEOUT))
	numBytes, _, err := us.conn.ReadFromUDP(us.packetBuffer)
	if err != nil {
		return nil, err
	}
	// Make a slice for the exact length of the packet and copy the packet
	// into it
	packet := make([]byte, numBytes)
	copy(packet, us.packetBuffer)
	return packet, nil
}

// processPackets is meant to be run in a goroutine and listens for packets
// until a shutdown signal is received
func (us *UDPSource) processPackets() {
	for {
		select {
		case <-us.shutdownChan:
			// Shutdown goroutine
			return
		default:
			packet, err := us.readPacket()
			if netError, ok := err.(net.Error); ok {
				// If timeout error, keep looping
				if !netError.Timeout() {
					// Panic if not a timeout error
					panic(err)
				}
			} else {
				us.outChan <- packet
			}
		}
	}
}

// listen sets up the UDPSource UDP connection on the given port
func (us *UDPSource) listen(port int) error {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	})
	if err != nil {
		return err
	}
	us.conn = conn
	go us.processPackets()
	return nil
}

// Packets returns a reference to the output packet channel
func (us *UDPSource) Packets() chan []byte {
	return us.outChan
}

// Close shuts down the goroutine and closes the connection
func (us *UDPSource) Close() {
	us.shutdownChan <- true
	us.conn.Close()
}

// SendPacketsAsUDP sends all the packets from the dataSource to the broadcast
// ip on the given port. Packets are spaced by the given delay duration.
func SendPacketsAsUDP(packetChan chan []byte, port int, delay time.Duration) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for packet := range packetChan {
		conn.Write(packet)
		time.Sleep(delay)
	}
}
