package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/gateway/helloworld/internal/server"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 8080, "port to listen on")
	mode = flag.String("mode", "dev", "mode to run")
)

func main() {
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	gsrv := server.NewGreeterGateway(context.Background(), fmt.Sprintf("localhost:%v", *port), opts)
	rsrv := server.NewGreeterRPCServer(*mode)

	log.Printf("http server and rpc server will run on the same port, server listening at http://localhost:%v\n", *port)
	log.Println("You can test it with: \n" + fmt.Sprintf(`    grpcurl -plaintext -d '{"name":"world"}' localhost:%v Greeter.SayHello`, *port))
	log.Println("You can test it with: \n" + fmt.Sprintf(`    curl -X POST http://localhost:%v/Greeter/SayHello -H "Content-Type: application/json" -d '{"name": "world"}'`, *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), server.MustServerMux(rsrv, gsrv)); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
