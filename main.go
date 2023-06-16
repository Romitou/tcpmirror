package main

import (
	"flag"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
)

func main() {
	interfaceName := flag.String("i", "eth0", "the interface to get packets from")
	inPort := flag.String("in", ":25565", "the port to listen on")
	outPort := flag.String("out", ":25566", "the port to mirror to")
	flag.Parse()

	handle, err := pcap.OpenLive(*interfaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	if err = handle.SetBPFFilter("tcp and port " + *inPort); err != nil {
		log.Fatal(err)
	}

	for packet := range gopacket.NewPacketSource(handle, handle.LinkType()).Packets() {
		var mirrorConn net.Conn
		mirrorConn, err = net.Dial("tcp", *outPort)
		if err != nil {
			log.Fatal(err)
		}

		_, err = mirrorConn.Write(packet.Data())
		if err != nil {
			log.Println(err)
		}
	}
}
