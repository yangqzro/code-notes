package server

import (
	"context"
	"fmt"
	"goexamples/gateway/helloworld/proto"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GreeterHTTPServer struct {
	mux *runtime.ServeMux
}

func (srv *GreeterHTTPServer) RawMux() *runtime.ServeMux {
	return srv.mux
}

func (srv *GreeterHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
}

func (srv *GreeterHTTPServer) Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), srv.mux)
}

// 创建一个反向代理服务器。用于将 RESTful http 请求转为 grpc 请求。
func MustNewGreeterHTTPServer(ctx context.Context, addr string, clientOpts []grpc.DialOption, opts ...runtime.ServeMuxOption) *GreeterHTTPServer {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandlerFromEndpoint(ctx, mux, addr, clientOpts); err != nil {
		panic(err)
	}
	return &GreeterHTTPServer{mux: mux}
}

func MustNewGreeterHTTPServerFromClient(ctx context.Context, client proto.GreeterClient, opts ...runtime.ServeMuxOption) *GreeterHTTPServer {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandlerClient(ctx, mux, client); err != nil {
		panic(err)
	}
	return &GreeterHTTPServer{mux: mux}
}

func MustNewGreeterHTTPServerFromConn(ctx context.Context, conn *grpc.ClientConn, opts ...runtime.ServeMuxOption) *GreeterHTTPServer {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandlerClient(ctx, mux, proto.NewGreeterClient(conn)); err != nil {
		panic(err)
	}
	return &GreeterHTTPServer{mux: mux}
}

func MustNewGreeterHTTPServerFromServer(ctx context.Context, server proto.GreeterServer, opts ...runtime.ServeMuxOption) *GreeterHTTPServer {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandlerServer(ctx, mux, server); err != nil {
		panic(err)
	}
	return &GreeterHTTPServer{mux: mux}
}
