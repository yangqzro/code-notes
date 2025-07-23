package client

import (
	"context"
	"goexamples/gateway/helloworld/proto"

	"google.golang.org/grpc"
)

type GreeterRPCClient struct {
	conn   *grpc.ClientConn
	client proto.GreeterClient
}

func (c *GreeterRPCClient) RawClient() proto.GreeterClient {
	return c.client
}

func (c *GreeterRPCClient) Close() error {
	return c.conn.Close()
}

func (cli *GreeterRPCClient) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	return cli.client.SayHello(ctx, in)
}

func NewGreeterRPCClient(addr string, opts ...grpc.DialOption) (*GreeterRPCClient, error) {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}
	return &GreeterRPCClient{conn: conn, client: proto.NewGreeterClient(conn)}, nil
}

func MustNewGreeterRPCClient(addr string, opts ...grpc.DialOption) *GreeterRPCClient {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		panic(err)
	}
	return &GreeterRPCClient{conn: conn, client: proto.NewGreeterClient(conn)}
}
