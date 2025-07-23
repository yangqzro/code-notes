# gRPC stream

一个更复杂的案例，演示了 gRPC 四种通信模式。

通信模式 | 方法名称
---|---
一元调用 (Unary) | GetPoem / GetPoemAll / UploadPoem / BatchUploadPoem
服务端流 (Server Stream) | GetPoemStream / GetPoemAllStream
客户端流 (Client Stream) | UploadPoemStream
双向流 (Bidirectional Stream) | BatchUploadPoemStream

```shell
cd grpc/examples/go/poem-stream

protoc --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative poem.proto # 编译 proto 文件

go run server/main.go # 1. 先运行服务端
go run client/main.go # 2. 再运行客户端
```
