package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/gateway/helloworld/internal/server"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	rpcPort  = flag.Int("rpc-port", 50051, "rpc server port to listen on")
	httpPort = flag.Int("http-port", 8080, "http server port to listen on")
	mode     = flag.String("mode", "dev", "mode to run")
)

func main() {
	flag.Parse()

	wait := sync.WaitGroup{}
	wait.Add(2)

	go func() {
		srv := server.NewGreeterRPCServer(*mode)
		listener := func(lis net.Listener) {
			log.Printf("rpc server listening at %v\n", lis.Addr())
			log.Println("You can test it with: \n" + fmt.Sprintf(`    grpcurl -plaintext -d '{"name":"world"}' localhost:%v Greeter.SayHello`, *rpcPort))
		}
		if err := srv.Listen(*rpcPort, listener); err != nil {
			log.Fatalf("failed to listen: %v\n", err)
		}
		wait.Done()
	}()

	go func() {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		srv := server.NewGreeterGateway(context.Background(), fmt.Sprintf("localhost:%v", *rpcPort), opts)
		log.Printf("http server listening at http://localhost:%v\n", *httpPort)
		log.Println("You can test it with: \n" + fmt.Sprintf(`    curl -X POST http://localhost:%v/Greeter/SayHello -H "Content-Type: application/json" -d '{"name": "world"}'`, *httpPort))
		if err := srv.Listen(*httpPort); err != nil {
			log.Fatalf("failed to listen: %v\n", err)
		}
		wait.Done()
	}()

	wait.Wait()
}
