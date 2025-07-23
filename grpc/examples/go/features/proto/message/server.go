package message

import (
	"context"
	"fmt"
	"goexamples/utils"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Server Handler 实现：
//
//  1. 接收客户端元数据 metadata
//  2. 依据情况，设置和发送 Header ？
//  3. 处理数据，接收客户端传入的数据，发送响应数据
//  4. 依据情况，设置 Trailer ？
//  5. 结束 rpc 响应
type MessageSrvServer struct {
	server *grpc.Server
	UnimplementedMessageServiceServer
}

func (s *MessageSrvServer) Listen(port int, onListen func(net.Listener)) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	if onListen != nil {
		onListen(lis)
	}
	return s.server.Serve(lis)
}

func (s *MessageSrvServer) Unary(ctx context.Context, in *Message) (*Message, error) {
	// metadata.FromIncomingContext 从上下文中获取客户端传来的元数据。
	// 服务端的响应元数据分为两种：Header 和 Trailer，Header 先于数据流发送，Trailer 在所有数据流结束后发送。
	// Header 用于传递‌前置处理信息‌，比如用户认证信息、会话数据、请求元信息，Trailer 携带‌后置处理结果，比如请求处理状态、统计数据、错误信息。
	// 服务端需要在处理 rpc 请求时显式调用 grpc.SendHeader / ServerStream.SenderHeader 发送 Header，如果没有发送，Header 会在首次返回响应数据时自动发送一次。
	// 在一次 rpc 生命周期内，Header 仅能发送一次，后续的发送将被忽略。
	// Trailer 无法通过显式调用函数发送，只有当所有响应数据流结束后由服务端自动发送。
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("server.Unary received metadata: %v\n", utils.String(md))
		// 在 Unary RPC 调用中，客户端获取响应 Header 不存在阻塞的情况，可以只设置不发送。
		grpc.SetHeader(ctx, GenerateServerMetadata("server.Unary header"))
		defer grpc.SetTrailer(ctx, GenerateServerMetadata("server.Unary trailer"))
	}

	log.Printf("server.Unary received message: %s\n", in.GetContent())
	return in, nil
}

func (s *MessageSrvServer) ClientStream(sin grpc.ClientStreamingServer[Message, MessageCollection]) error {
	if md, ok := metadata.FromIncomingContext(sin.Context()); ok {
		log.Printf("server.ClientStream received metadata: %v\n", utils.String(md))
		// 由于 ServerStream 会在首次发送响应数据时自动发送 Header，因此必须确保 ServerStream.SendHeader() 的调用先于任何 ServerStream.Send() 操作，否则后续调用 ServerStream.SendHeader() 将无效。
		// ClientStream.Header() 是个阻塞方法。
		// 在流式 rpc 通信中，如果客户端 Header 读取（ClientStream.Header()）先于响应数据读取（ClientStream.Recv()、ClientStreamingClient.CloseAndRecv()），
		// 则服务端必须严格确保在调用 ServerStream.Send() 发送数据前完成 Header 的发送，否则会导致双向阻塞死锁。
		sin.SendHeader(GenerateServerMetadata("server.ClientStream header"))
		defer sin.SetTrailer(GenerateServerMetadata("server.ClientStream trailer"))
	}
	mc := []*Message{}
	for {
		in, err := sin.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("server.ClientStream received message: %s\n", in.GetContent())
		mc = append(mc, in)
	}
	return sin.SendAndClose(&MessageCollection{Value: mc})
}

func (s *MessageSrvServer) ServerStream(in *MessageCollection, sout grpc.ServerStreamingServer[Message]) error {
	if md, ok := metadata.FromIncomingContext(sout.Context()); ok {
		log.Printf("server.ServerStream received metadata: %v\n", utils.String(md))
		// 由于 ServerStream 会在首次发送响应数据时自动发送 Header，因此必须确保 ServerStream.SendHeader() 的调用先于任何 ServerStream.Send() 操作，否则后续调用 ServerStream.SendHeader() 将无效。
		sout.SendHeader(GenerateServerMetadata("server.ServerStream header"))
		defer sout.SetTrailer(GenerateServerMetadata("server.ServerStream trailer"))
	}
	for _, m := range in.GetValue() {
		log.Printf("server.ServerStream received message: %s\n", m.GetContent())
		if err := sout.Send(m); err != nil {
			return err
		}
	}
	return nil
}

func (s *MessageSrvServer) BidirectionalStream(sbin grpc.BidiStreamingServer[Message, Message]) error {
	if md, ok := metadata.FromIncomingContext(sbin.Context()); ok {
		log.Printf("server.BidirectionalStream received metadata: %v\n", utils.String(md))
		// 由于 ServerStream 会在首次发送响应数据时自动发送 Header，因此必须确保 ServerStream.SendHeader() 的调用先于任何 ServerStream.Send() 操作，否则后续调用 ServerStream.SendHeader() 将无效。
		sbin.SendHeader(GenerateServerMetadata("server.BidirectionalStream header"))
		defer sbin.SetTrailer(GenerateServerMetadata("server.BidirectionalStream trailer"))
	}
	for {
		in, err := sbin.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("server.BidirectionalStream received message: %s\n", in.GetContent())
		if err := sbin.Send(in); err != nil {
			return err
		}
	}
	return nil
}

func NewMessageSrvServer(opts ...grpc.ServerOption) *MessageSrvServer {
	server := grpc.NewServer(opts...)
	srv := &MessageSrvServer{server: server}
	RegisterMessageServiceServer(server, srv)
	return srv
}
