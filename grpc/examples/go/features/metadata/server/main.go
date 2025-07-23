package main

import (
	"flag"
	"goexamples/features/proto/message"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "port to listen on")
)

func main() {
	flag.Parse()
	server := message.NewMessageSrvServer()
	onListen := func(lis net.Listener) {
		log.Printf("server listening at %v\n", lis.Addr())
	}
	if err := server.Listen(*port, onListen); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
