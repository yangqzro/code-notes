#!/usr/bin/env bash
protoc --go_out=. --go-grpc_out=. ./message.proto
