package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/gateway/helloworld/internal/server"
	"goexamples/utils"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerMux struct {
	rpcServer  http.Handler
	httpServer http.Handler
}

func (srv *ServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request: url=%v, method=%v, header=%s", r.URL, r.Method, utils.String(r.Header))
	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		srv.rpcServer.ServeHTTP(w, r) // 处理 gRPC 请求
	} else {
		srv.httpServer.ServeHTTP(w, r)
	}
}

func MustNewServerMux(rpcServer, httpServer http.Handler) *ServerMux {
	if rpcServer == nil || httpServer == nil {
		panic("rpcServer or httpServer is nil")
	}
	return &ServerMux{rpcServer: rpcServer, httpServer: httpServer}
}

func EnableH2C(handler http.Handler) http.Handler {
	return h2c.NewHandler(handler, &http2.Server{})
}

var (
	port = flag.Int("port", 8080, "port to listen on")
	mode = flag.String("mode", "dev", "mode to run")
)

func main() {
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	hsrv := server.MustNewGreeterHTTPServer(context.Background(), fmt.Sprintf("localhost:%v", *port), opts)
	rsrv := server.NewGreeterRPCServer(*mode)

	mux := MustNewServerMux(rsrv, hsrv)
	log.Printf("http server and rpc server will run on the same port, server listening at http://localhost:%v\n", *port)

	log.Println("You can test it with: \n" + fmt.Sprintf(`    grpcurl -plaintext -d '{"name":"world"}' localhost:%v Greeter.SayHello`, *port))
	log.Println("You can test it with: \n" + fmt.Sprintf(`    curl -X POST http://localhost:%v/Greeter/SayHello -H "Content-Type: application/json" -d '{"name": "world"}'`, *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), EnableH2C(mux)); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
