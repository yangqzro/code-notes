# gRPC gateway

```shell
cd grpc/examples/go/gateway

# 编译 proto 文件
buf generate
# 或者
protoc \
  --go_out=proto \
  --go_opt=paths=source_relative \
  --go-grpc_out=proto \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=proto \
  --grpc-gateway_opt=paths=source_relative,generate_unbound_methods=true \
helloworld.proto

#  运行测试 helloworld
go run cmd/helloworld/main.go
grpcurl -plaintext -d '{"name": "world"}' localhost:50051 Greeter.SayHello
curl -X POST http://localhost:8080/Greeter/SayHello -H "Content-Type: application/json" -d '{"name": "world"}'

#  运行测试 sameport
go run cmd/sameport/main.go
grpcurl -plaintext -d '{"name": "world"}' localhost:8080 Greeter.SayHello
curl -X POST http://localhost:8080/Greeter/SayHello -H "Content-Type: application/json" -d '{"name": "world"}'
```
