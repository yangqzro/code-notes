package server

import (
	"context"
	"fmt"
	"goexamples/gateway/helloworld/proto"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GreeterRPCServer struct {
	proto.UnimplementedGreeterServer
	server *grpc.Server
}

func (srv *GreeterRPCServer) RawServer() *grpc.Server {
	return srv.server
}

func (srv *GreeterRPCServer) GetServiceInfo() map[string]grpc.ServiceInfo {
	return srv.server.GetServiceInfo()
}

func (srv *GreeterRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.server.ServeHTTP(w, r)
}

func (srv *GreeterRPCServer) Listen(port int, listener func(net.Listener)) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	if listener != nil {
		listener(lis)
	}
	return srv.server.Serve(lis)
}

func (srv *GreeterRPCServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hello " + req.Name}, nil
}

func NewGreeterRPCServer(mode string, opts ...grpc.ServerOption) *GreeterRPCServer {
	server := grpc.NewServer(opts...)
	srv := &GreeterRPCServer{server: server}
	if mode == "dev" {
		reflection.Register(server)
	}
	proto.RegisterGreeterServer(server, srv)
	return srv
}
