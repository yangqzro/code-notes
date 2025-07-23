# gRPC hello world

一个 gRPC 客户端和服务端的简单例子。

```shell
cd grpc/examples/go/helloworld

# 编译 proto 文件
protoc --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative helloworld.proto

go run server/main.go # 1. 先运行服务端
go run client/main.go # 2. 再运行客户端
```
