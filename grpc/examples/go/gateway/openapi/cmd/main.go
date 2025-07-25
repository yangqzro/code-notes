package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/gateway/openapi/internal/client"
	"goexamples/gateway/openapi/internal/model"
	"goexamples/gateway/openapi/internal/server"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 8080, "port to listen on")
	mode = flag.String("mode", "dev", "mode to run")
)

func main() {
	flag.Parse()

	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get current directory: %v\n", err)
	}

	rsrv := server.NewUserRPCServer(*mode)
	rsrv.SetModel(model.NewUserModel().MustLoad(filepath.Join(root, "testdata", "users.json")))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(client.PostAutoFillFieldMask)}
	gsrv := server.NewUserGateway(context.Background(), fmt.Sprintf("localhost:%v", *port), opts, runtime.WithMiddlewares(server.BodyBufferMiddleware))

	fsrv := server.StaticServer(filepath.Join(root, "third_party", "openapi"))

	mux := server.NewServerMux(
		rsrv, func(r *http.Request) bool {
			return r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc")
		},
		gsrv, func(r *http.Request) bool {
			return strings.HasPrefix(r.URL.Path, "/api")
		},
		fsrv, nil,
	)
	log.Printf("http server and rpc server will run on the same port, server listening at http://localhost:%v, you can visit it to get api docs\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), server.EnableH2C(mux)); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
