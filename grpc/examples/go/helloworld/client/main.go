package main

import (
	"context"
	"flag"
	"goexamples/helloworld/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.GreeterClient
}

func (c *Client) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return c.client.SayHello(ctx, in)
}

func (c *Client) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return c.client.SayHelloAgain(ctx, in)
}

func (c *Client) Close() {
	c.conn.Close()
}

func NewClient(addr string) *Client {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &Client{
		conn:   conn,
		client: pb.NewGreeterClient(conn),
	}
}

var (
	addr = flag.String("addr", "localhost:50051", "port to connect to")
	name = flag.String("name", "world", "name to greet")
)

func main() {
	flag.Parse()

	c := NewClient(*addr)
	defer c.Close()

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	r, err = c.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
