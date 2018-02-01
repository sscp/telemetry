package collector

import (
	"context"
	"net"
	"time"
)

// UDPListenTimeout is the time to wait for the next packet
const udpListenTimeout = 100 * time.Millisecond

type UDPPacketSource struct {
	outChan      chan *ContextPacket
	doneChan     chan bool
	conn         *net.UDPConn
	packetBuffer []byte
}

func NewUDPPacketSource(port int) (PacketSource, error) {
	// Listen to the zero port for IPv4 to catch any packet to that port
	// This will catch broadcast packets from the car
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	})
	if err != nil {
		return nil, err
	}
	return &UDPPacketSource{
		outChan:      make(chan *ContextPacket),
		doneChan:     make(chan bool),
		conn:         conn,
		packetBuffer: make([]byte, 1000), // Max packet size is 512
	}, nil
}

// Packets is the stream of packets received from UDP
// It is simply a reference to outChan
func (ups *UDPPacketSource) Packets() <-chan *ContextPacket {
	return ups.outChan
}

func (ups *UDPPacketSource) Listen() {

	go func() {
		for {
			select {
			case <-ups.doneChan:
				// Close Conn and shutdown goroutine
				ups.conn.Close()
				return
			default:
				packet, err := ups.readPacket()
				if netError, ok := err.(net.Error); ok {
					// If timeout error, keep looping
					if !netError.Timeout() {
						// Panic if not a timeout error
						panic(err)
					}
				} else {
					ctx := context.TODO()
					ups.outChan <- &ContextPacket{
						ctx:    ctx,
						packet: packet,
					}
				}
			}
		}
	}()

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

func (ups *UDPPacketSource) Close() {
	ups.doneChan <- true
	close(ups.doneChan)
	close(ups.outChan)
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
