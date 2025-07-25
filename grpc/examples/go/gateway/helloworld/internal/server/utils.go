package server

import (
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func EnableH2C(handler http.Handler) http.Handler {
	return h2c.NewHandler(handler, &http2.Server{})
}

func MustServerMux(rpcServer, httpServer http.Handler) http.Handler {
	if rpcServer == nil || httpServer == nil {
		panic("rpcServer or httpServer is nil")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			rpcServer.ServeHTTP(w, r) // 处理 gRPC 请求
		} else {
			httpServer.ServeHTTP(w, r)
		}
	})
}
