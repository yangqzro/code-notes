package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/helloworld/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	listener net.Listener
	server   *grpc.Server
	proto.UnimplementedGreeterServer
}

func (*Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (*Server) SayHelloAgain(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "Hello " + in.GetName() + " again"}, nil
}

func NewServer(port int) *Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	srv := &Server{server: gs, listener: lis}
	proto.RegisterGreeterServer(gs, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return srv
}

var (
	port = flag.Int("port", 50051, "port to listen on")
)

func main() {
	flag.Parse()
	_ = NewServer(*port)
}
