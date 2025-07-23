#!/usr/bin/env bash
# https://protobuf.dev/installation/
brew install protobuf

# https://grpc.io/docs/languages/go/quickstart/
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
