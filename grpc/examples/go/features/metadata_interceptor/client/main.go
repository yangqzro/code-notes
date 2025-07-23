package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/features/pb/message"
	"goexamples/utils"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func unaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 追加元数据
	ctx = metadata.AppendToOutgoingContext(ctx, "token", utils.RandString(8))
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		log.Printf("client unary interceptor called: call metadata.FromOutgoingContext=%v\n", utils.String(md))
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// 追加元数据
	ctx = metadata.AppendToOutgoingContext(ctx, "token", utils.RandString(8))
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		log.Printf("client stream interceptor called: call metadata.FromOutgoingContext=%v\n", utils.String(md))
	}
	return streamer(ctx, desc, cc, method, opts...)
}

var (
	addr = flag.String("addr", "localhost:50051", "addr to connect to")
)

func main() {
	flag.Parse()
	client, err := message.NewMessageSrvClient(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryInterceptor),
		grpc.WithStreamInterceptor(streamInterceptor),
	)
	if err != nil {
		log.Fatalf("did not start client: %v\n", err)
	}
	defer client.Close()

	func() {
		m := "hello world"
		log.Printf("main.client.Unary send message: %s\n", m)

		// header 和 trailer 用于接收服务端返回的元数据
		var header, trailer metadata.MD
		if out, err := client.Unary(context.Background(), &message.Message{Content: m}, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
			log.Fatalf("main.client.Unary failed: %v\n", err)
		} else {
			log.Printf("main.client.Unary header: %v\n", utils.String(header))
			log.Printf("main.client.Unary trailer: %v\n", utils.String(trailer))
			log.Printf("main.client.Unary response: %v\n", out.GetContent())
		}
	}()

	func() {
		mc := []*message.Message{{Content: "client stream message 1"}, {Content: "client stream message 2"}, {Content: "client stream message 3"}}
		fmt.Println()
		log.Println("main.client.ClientStream send message")

		var header, trailer metadata.MD
		if out, err := client.ClientStream(context.Background(), mc, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
			log.Fatalf("main.client.ClientStream failed: %v\n", err)
		} else {
			log.Printf("main.client.ClientStream header: %v\n", utils.String(header))
			log.Printf("main.client.ClientStream trailer: %v\n", utils.String(trailer))
			for _, m := range out {
				log.Printf("main.client.ClientStream response: %v\n", m.GetContent())
			}
		}
	}()

	func() {
		mc := []*message.Message{{Content: "server stream message 1"}, {Content: "server stream message 2"}, {Content: "server stream message 3"}}
		fmt.Println()
		log.Println("main.client.ServerStream send message")

		var header, trailer metadata.MD
		if out, err := client.ServerStream(context.Background(), mc, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
			log.Fatalf("main.client.ServerStream failed: %v\n", err)
		} else {
			log.Printf("main.client.ServerStream header: %v\n", utils.String(header))
			log.Printf("main.client.ServerStream trailer: %v\n", utils.String(trailer))
			for _, m := range out {
				log.Printf("main.client.ServerStream response: %v\n", m.GetContent())
			}
		}
	}()

	func() {
		mc := []*message.Message{{Content: "bidirectional stream message 1"}, {Content: "bidirectional stream message 2"}, {Content: "bidirectional stream message 3"}}
		fmt.Println()
		log.Println("main.client.BidirectionalStream send message")

		var header, trailer metadata.MD
		if out, err := client.BidirectionalStream(context.Background(), mc, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
			log.Fatalf("main.client.BidirectionalStream failed: %v\n", err)
		} else {
			log.Printf("main.client.BidirectionalStream header: %v\n", utils.String(header))
			log.Printf("main.client.BidirectionalStream trailer: %v\n", utils.String(trailer))
			for _, m := range out {
				log.Printf("main.client.BidirectionalStream response: %v\n", m.GetContent())
			}
		}
	}()
}
