package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/mdlayher/vsock"
)

var (
	destHostPort string
	port         uint32
)

func init() {
	flag.StringVar(&destHostPort, "h", "localhost:22", "Destination host/port")

	port64 := flag.Uint64("l", 22, "Port to listen on")
	port = uint32(*port64)
	flag.Parse()
}

func main() {
	ln, err := vsock.Listen(port)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%s\n", err)
			continue
		}
		defer conn.Close()

		go func(conn net.Conn) {
			lconn, err := net.Dial("tcp", destHostPort)
			if err != nil {
				panic(err)
			}
			defer lconn.Close()

			go func() {
				io.Copy(lconn, conn)
			}()
			io.Copy(conn, lconn)
		}(conn)
	}
}
