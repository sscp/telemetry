package collector

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestUDPSendRecv(t *testing.T) {
	src, err := NewUDPPacketSource(3000)
	if err != nil {
		t.Errorf("Error creating packet source: %v", err)
	}

	packetChan := make(chan []byte)
	// Send all the packets in the channel
	go SendPacketsAsUDP(packetChan, 3000, time.Duration(0))

	// Listen for those same packets
	src.Listen()
	defer src.Close()

	for i := 1; i < 512; i++ {
		// Make a random packet
		packet := make([]byte, i)
		rand.Read(packet)

		// Packet sent to the send channel
		packetChan <- packet

		// Listen for the packet on the recv channel
		outPacket := <-src.Packets()

		// Check that everything made it
		if !reflect.DeepEqual(packet, outPacket.packet) {
			t.Errorf("Output packet, %v, does not match input packet %v", outPacket.packet, packet)
		}
	}
}
