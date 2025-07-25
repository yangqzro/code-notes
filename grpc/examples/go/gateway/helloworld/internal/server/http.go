package server

import (
	"context"
	"fmt"
	"goexamples/gateway/helloworld/proto"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GreeterGateway struct {
	mux *runtime.ServeMux
}

func (srv *GreeterGateway) RawMux() *runtime.ServeMux {
	return srv.mux
}

func (srv *GreeterGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
}

func (srv *GreeterGateway) Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), srv.mux)
}

// 创建一个反向代理服务器。用于将 RESTful http 请求转为 grpc 请求。
func NewGreeterGateway(ctx context.Context, addr string, clientOpts []grpc.DialOption, opts ...runtime.ServeMuxOption) *GreeterGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandlerFromEndpoint(ctx, mux, addr, clientOpts); err != nil {
		panic(err)
	}
	return &GreeterGateway{mux: mux}
}

func NewGreeterGatewayFromConn(ctx context.Context, conn *grpc.ClientConn, opts ...runtime.ServeMuxOption) *GreeterGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandler(ctx, mux, conn); err != nil {
		panic(err)
	}
	return &GreeterGateway{mux: mux}
}

func NewGreeterGatewayFromClient(ctx context.Context, client proto.GreeterClient, opts ...runtime.ServeMuxOption) *GreeterGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandlerClient(ctx, mux, client); err != nil {
		panic(err)
	}
	return &GreeterGateway{mux: mux}
}

func NewGreeterGatewayFromServer(ctx context.Context, server proto.GreeterServer, opts ...runtime.ServeMuxOption) *GreeterGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterGreeterHandlerServer(ctx, mux, server); err != nil {
		panic(err)
	}
	return &GreeterGateway{mux: mux}
}
