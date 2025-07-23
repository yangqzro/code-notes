package main

import (
	"context"
	"flag"
	"goexamples/features/pb/message"
	"goexamples/utils"
	"log"
	"net"

	"google.golang.org/grpc"
)

// req: 客户端传来的请求参数
// info.FullName: 客户端要调用的完整 RPC 方法名称，格式：/package.Service/Method
// info.Server: 用户注册的服务端实现，本例中指的 main 方法的 server 变量
// handler: 实际处理 rpc 请求的方法
func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("server unary interceptor called: req=%v, info=%v", utils.String(req), utils.String(info))
	return handler(ctx, req)
}

func unaryInterceptor2(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("server unary2 interceptor called: req=%v, info=%v", utils.String(req), utils.String(info))
	return handler(ctx, req)
}

// srv: 用户注册的服务端实现，本例中指的 main 方法的 server 变量
// ss: 服务端流对象，实际类型会根据不同的流模式变化，可以通过类型断言获取具体流对象，本例中
//
//	客户端流 (ClientStream) -> *grpc.GenericServerStream[Message, MessageCollection]
//	服务端流 (ServerStream) -> *grpc.GenericServerStream[MessageCollection, Message]
//	双向流 (BidiStream) -> *grpc.GenericServerStream[Message, Message]
//
// info.FullName: 客户端要调用的完整 RPC 方法名称，格式：/package.Service/Method
// info.IsClientStream: 判断 rpc 是不是客户端流模式，双向流时为 true
// info.IsServerStream: 判断 rpc 是不是服务端流模式，双向流时为 true
// handler: 实际处理 rpc 请求的方法
func streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("server stream interceptor called: srv=%v, ss=%T, info=%v\n", utils.String(srv), utils.String(ss), utils.String(info))
	return handler(srv, ss)
}

var (
	port = flag.Int("port", 50051, "port to listen on")
)

func main() {
	flag.Parse()
	server := message.NewMessageSrvServer(
		// 配置一个一元拦截器，如果需要多个拦截器，请使用 grpc.ChainUnaryInterceptor
		grpc.UnaryInterceptor(unaryInterceptor),
		// 配置一个流拦截器，如果需要多个拦截器，请使用 grpc.ChainStreamInterceptor
		grpc.StreamInterceptor(streamInterceptor),
		// 多个拦截器按注册的顺序执行
		// grpc.ChainUnaryInterceptor(unaryInterceptor, unaryInterceptor2),
	)
	onListen := func(lis net.Listener) {
		log.Printf("server listening at %v\n", lis.Addr())
	}
	if err := server.Listen(*port, onListen); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
