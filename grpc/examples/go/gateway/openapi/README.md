# gRPC gateway

```shell
cd grpc/examples/go/gateway/openapi

# 编译 proto 文件
buf generate

#  运行测试 openapi
go run cmd/main.go
grpcurl -plaintext -d '{"user": {"email": "zhangsan@example.com", "name": "zhangsan"}}' localhost:8080 user.UserService.AddUser
curl -X POST http://localhost:8080/api/v1/users -H "Content-Type: application/json" -d '{"email": "zhangsan@example.com", "name": "zhangsan"}'
curl -X DELETE http://localhost:8080/api/v1/users/1
curl -X PATCH http://localhost:8080/api/v1/users/1 -H "Content-Type: application/json" -d '{"email": "lisi@example.com"}'
curl -X GET http://localhost:8080/api/v1/users/1
curl -X GET http://localhost:8080/api/v1/users
```

[https://github.com/johanbrandhorst/grpc-gateway-boilerplate](https://github.com/johanbrandhorst/grpc-gateway-boilerplate)
