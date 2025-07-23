package main

import (
	"context"
	"errors"
	"flag"
	"goexamples/features/proto/message"
	"goexamples/utils"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	ErrMetadataMiss = errors.New("metadata miss")
	ErrTokenMiss    = errors.New("token miss")
)

func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMetadataMiss
	}
	log.Printf("server unary interceptor called: md=%v", utils.String(md))

	if md.Get("token") != nil {
		grpc.SetHeader(ctx, metadata.Pairs("user", utils.RandString(4)))
		grpc.SetTrailer(ctx, metadata.Pairs("user", utils.RandString(4)))
	} else {
		return nil, ErrTokenMiss
	}
	return handler(ctx, req)
}

func streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return ErrMetadataMiss
	}
	log.Printf("server stream interceptor called: md=%v", utils.String(md))

	if md.Get("token") != nil {
		ss.SetHeader(metadata.Pairs("user", utils.RandString(4)))
		ss.SetTrailer(metadata.Pairs("user", utils.RandString(4)))
	} else {
		return ErrTokenMiss
	}
	return handler(srv, ss)
}

var (
	port = flag.Int("port", 50051, "port to listen on")
)

func main() {
	flag.Parse()
	server := message.NewMessageSrvServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)
	onListen := func(lis net.Listener) {
		log.Printf("server listening at %v\n", lis.Addr())
	}
	if err := server.Listen(*port, onListen); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
