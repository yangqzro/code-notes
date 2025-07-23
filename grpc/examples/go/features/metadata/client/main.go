package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/features/proto/message"
	"goexamples/utils"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	addr = flag.String("addr", "localhost:50051", "addr to connect to")
)

func main() {
	flag.Parse()
	client, err := message.NewMessageSrvClient(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not start client: %v\n", err)
	}
	defer client.Close()

	// metadata 是 gRPC 中用于传递元数据的机制，类似于 HTTP 的 Header。它的生命周期就是一次 rpc 调用。
	// metadata.MD 是一个 map[string][]string 的别名，意味着同一个键可以对应多个值。
	// metadata.New 由于 map 的特性，初始化不能声明重复的键。如果需要同一个键对应多个值，可以使用 metadata.Pairs。
	// metadata.New 和 metadata.Pairs 在初始化时，键会自动转为小写，键名只支持数字、大小写字母、下划线（_）和连字符（-）。
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.DateTime),
		// 同一个键的多个值合并成一个切片 metadata.MD{"timestamp": [time.Now().Format(time.DateTime), time.Now().Format(time.RFC3339)]}。
		// "timestamp", time.Now().Format(time.RFC3339),
		// 以 grpc- 开头的键仅供 grpc 内部使用，如果在元数据中设置可能会导致错误，所以建名尽量不要以 grpc- 开头。
		// "grpc-timestamp", "example",
	)

	func() {
		m := "hello world"
		log.Printf("main.client.Unary send message: %s\n", m)

		// metadata.NewOutgoingContext 创建一个新的上下文，并将元数据附加到上下文中。如果上下文已经有了元数据，它会覆盖已有的元数据。
		// 如果想要追加元数据，可以使用 metadata.AppendToOutgoingContext 函数。
		// 该函数将传入的元数据与已有的元数据合并，并且返回一个新的上下文，在拦截器中比较实用。
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		// header 和 trailer 用于接收服务端返回的元数据
		var header, trailer metadata.MD
		if out, err := client.Unary(ctx, &message.Message{Content: m}, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
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

		ctx := metadata.NewOutgoingContext(context.Background(), md)
		var header, trailer metadata.MD
		if out, err := client.ClientStream(ctx, mc, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
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

		ctx := metadata.NewOutgoingContext(context.Background(), md)
		var header, trailer metadata.MD
		if out, err := client.ServerStream(ctx, mc, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
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

		ctx := metadata.NewOutgoingContext(context.Background(), md)
		var header, trailer metadata.MD
		if out, err := client.BidirectionalStream(ctx, mc, grpc.Header(&header), grpc.Trailer(&trailer)); err != nil {
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
