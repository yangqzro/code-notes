package server

import (
	"bytes"
	"context"
	"fmt"
	"goexamples/gateway/openapi/proto"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var (
	BodyReaderContext = &struct{}{}
)

type UserGateway struct {
	mux *runtime.ServeMux
}

func (srv *UserGateway) RawMux() *runtime.ServeMux {
	return srv.mux
}

func (srv *UserGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
}

func (srv *UserGateway) Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), srv.mux)
}

// 创建一个反向代理服务器。用于将 RESTful http 请求转为 grpc 请求。
func NewUserGateway(ctx context.Context, addr string, clientOpts []grpc.DialOption, opts ...runtime.ServeMuxOption) *UserGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addr, clientOpts); err != nil {
		panic(err)
	}
	return &UserGateway{mux: mux}
}

func NewUserGatewayFromConn(ctx context.Context, conn *grpc.ClientConn, opts ...runtime.ServeMuxOption) *UserGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterUserServiceHandler(ctx, mux, conn); err != nil {
		panic(err)
	}
	return &UserGateway{mux: mux}
}

func NewUserGatewayFromClient(ctx context.Context, client proto.UserServiceClient, opts ...runtime.ServeMuxOption) *UserGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterUserServiceHandlerClient(ctx, mux, client); err != nil {
		panic(err)
	}
	return &UserGateway{mux: mux}
}

func NewUserGatewayFromServer(ctx context.Context, server proto.UserServiceServer, opts ...runtime.ServeMuxOption) *UserGateway {
	mux := runtime.NewServeMux(opts...)
	if err := proto.RegisterUserServiceHandlerServer(ctx, mux, server); err != nil {
		panic(err)
	}
	return &UserGateway{mux: mux}
}

func BodyBufferMiddleware(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		if r.Method == http.MethodPost || r.Method == http.MethodPatch || r.Method == http.MethodPut {
			// 由于 io.ReadCloser 的限制，body 只能读取一次
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// 这里要重新设置 body，否则会导致后续的读取失败
			r.Body = io.NopCloser(bytes.NewReader(body))
			next(w, r.WithContext(context.WithValue(r.Context(), BodyReaderContext, bytes.NewReader(body))), pathParams)
		} else {
			next(w, r, pathParams)
		}
	}
}
