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
)

// method: 客户端要调用的完整 rpc 方法名称，格式：/package.Service/Method
// req: 请求参数
// reply: 响应参数，指针类型，调用 invoker 后，会将响应参数赋值给 reply，即调用 invoker 后 reply 才有值
// cc: 客户端连接实例，本例中是 main 方法中 client.GetConn() 方法返回的连接实例
// invoker: 实际执行 rpc 调用的处理器
// opts: 本次 rpc 调用携带的配置选项
func unaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("client unary interceptor called: method=%v, req=%v, reply=%v\n", utils.String(method), utils.String(req), utils.String(reply))
	return err
}

func unaryInterceptor2(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("client unary2 interceptor called: method=%v, req=%v, reply=%v\n", utils.String(method), utils.String(req), utils.String(reply))
	return err
}

// desc: 流描述信息，包含流的类型（客户端流、服务器流、双向流）和消息类型
// desc.StreamName: 客户端要调用 rpc 方法名称，不包含服务名
// desc.Handler: 该函数由 .proto 文件中 service 定义的 gRPC 方法自动生成，命名格式为 _Service_Method_Handler。客户端无需关注此字段，它仅用于服务端内部注册处理gRPC请求的业务逻辑函数。在此例中
//
//	ClientStream -> _MessageService_ClientStream_Handler
//	ServerStream -> _MessageService_ServerStream_Handler
//	BidirectionalStream -> _MessageService_BidirectionalStream_Handler
//
// desc.ServerStreams: 判断 rpc 是不是客户端流模式，双向流时为 true，表示服务端可推送流式响应
// desc.ClientStreams: 判断 rpc 是不是服务端流模式，双向流时为 true，表示客户端可推送流式请求
// cc: 客户端连接实例，本例中是 main 方法中 client.GetConn() 方法返回的连接实例
// method: 客户端要调用的完整 rpc 方法名称，格式：/package.Service/Method
// streamer: 实际创建流的处理器
// opts: 本次 rpc 调用携带的配置选项
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("client stream interceptor called: desc=%v, method=%v\n", utils.String(desc), utils.String(method))
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
		// 配置一个一元拦截器，如果需要多个拦截器，请使用 grpc.WithChainUnaryInterceptor
		grpc.WithUnaryInterceptor(unaryInterceptor),
		// 配置一个流拦截器，如果需要多个拦截器，请使用 grpc.WithChainStreamInterceptor
		grpc.WithStreamInterceptor(streamInterceptor),
		// 多个拦截器按注册的顺序执行
		// grpc.WithChainUnaryInterceptor(unaryInterceptor, unaryInterceptor2),
	)
	if err != nil {
		log.Fatalf("did not start client: %v\n", err)
	}
	defer client.Close()

	func() {
		m := "hello world"
		log.Printf("main.client.Unary send message: %s\n", m)

		if out, err := client.Unary(context.Background(), &message.Message{Content: m}); err != nil {
			log.Fatalf("main.client.Unary failed: %v\n", err)
		} else {
			log.Printf("main.client.Unary response: %v\n", out.GetContent())
		}
	}()

	func() {
		mc := []*message.Message{{Content: "client stream message 1"}, {Content: "client stream message 2"}, {Content: "client stream message 3"}}
		fmt.Println()
		log.Println("main.client.ClientStream send message")

		if out, err := client.ClientStream(context.Background(), mc); err != nil {
			log.Fatalf("main.client.ClientStream failed: %v\n", err)
		} else {
			for _, m := range out {
				log.Printf("main.client.ClientStream response: %v\n", m.GetContent())
			}
		}
	}()

	func() {
		mc := []*message.Message{{Content: "server stream message 1"}, {Content: "server stream message 2"}, {Content: "server stream message 3"}}
		fmt.Println()
		log.Println("main.client.ServerStream send message")

		if out, err := client.ServerStream(context.Background(), mc); err != nil {
			log.Fatalf("main.client.ServerStream failed: %v\n", err)
		} else {
			for _, m := range out {
				log.Printf("main.client.ServerStream response: %v\n", m.GetContent())
			}
		}
	}()

	func() {
		mc := []*message.Message{{Content: "bidirectional stream message 1"}, {Content: "bidirectional stream message 2"}, {Content: "bidirectional stream message 3"}}
		fmt.Println()
		log.Println("main.client.BidirectionalStream send message")

		if out, err := client.BidirectionalStream(context.Background(), mc); err != nil {
			log.Fatalf("main.client.BidirectionalStream failed: %v\n", err)
		} else {
			for _, m := range out {
				log.Printf("main.client.BidirectionalStream response: %v\n", m.GetContent())
			}
		}
	}()
}
