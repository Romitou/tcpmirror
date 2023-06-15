package main

import (
	"flag"
	"io"
	"log"
	"net"
)

func main() {
	inPort := flag.String("in", ":25565", "the port to listen on")
	outPort := flag.String("out", ":25566", "the port to mirror to")
	flag.Parse()

	listeningPort, err := net.Listen("tcp", *inPort)
	if err != nil {
		log.Fatal(err)
	}

	defer func(listeningPort net.Listener) {
		err = listeningPort.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listeningPort)

	mirrorConn, err := net.Dial("tcp", *outPort)
	if err != nil {
		log.Fatal(err)
	}

	for {
		incomingConn, _ := listeningPort.Accept()
		go func() {
			_, err = io.Copy(incomingConn, mirrorConn)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}
