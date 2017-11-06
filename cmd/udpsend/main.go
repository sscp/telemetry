package main

import (
	"flag"
	"fmt"
	"github.com/sscp/naturallight-telemetry/datasources"
	"os"
	"time"
)

func main() {
	delayStr := flag.String("delay", "0ms", "delay between sending packets")
	port := flag.Int("port", 33333, "port to send packets on")
	flag.Parse()

	filename := flag.Arg(0)

	delay, err := time.ParseDuration(*delayStr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sending packets from %v on port %v with delay %v", filename, *port, delay)

	sendPacketsFromBlog(filename, *port, delay)
}

func sendPacketsFromBlog(filename string, port int, delay time.Duration) {
	blog, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	ds := datasources.ReadPackets(blog, 0)

	datasources.SendPacketsAsUDP(ds.Packets(), port, delay)
}
