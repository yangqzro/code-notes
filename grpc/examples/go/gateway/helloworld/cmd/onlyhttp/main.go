package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/gateway/helloworld/internal/server"
	"log"
)

var (
	port = flag.Int("port", 8080, "port to listen on")
)

func main() {
	flag.Parse()

	// 通过 Server 业务实现创建 http 服务器。该服务器仅对外暴露业务接口，由于未初始化 gRPC 服务实例，外部客户端无法建立 gRPC 协议连接进行远程调用。
	// http 服务器只负责将 JSON 数据转发给业务方法，中间不涉及 protobuf 的编解码。
	hsrv := server.MustNewGreeterHTTPServerFromServer(context.Background(), &server.GreeterRPCServer{})
	log.Printf("http server listening at http://localhost:%v\n", *port)
	log.Println("rpc server is not started")

	log.Println("You can test it with: \n" + fmt.Sprintf(`    curl -X POST http://localhost:%v/Greeter/SayHello -H "Content-Type: application/json" -d '{"name": "world"}'`, *port))
	if err := hsrv.Listen(*port); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
