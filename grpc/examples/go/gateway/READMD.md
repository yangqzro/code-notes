# grpc-gateway

gRPC-Gateway 能够根据 protobuf 服务定义，自动生成一个反向代理服务器，将 RESTful HTTP 请求转换为 gRPC 请求。

**它本质上是一个 gRPC 客户端。**作为客户端，它的核心功能是将接收到的 JSON 请求数据构造成 gRPC 请求格式，通过内置的 gRPC 客户端发送给服务端，并将服务端返回的响应数据重新封装为 HTTP 响应返回给客户端。

- 路由绑定，将 URL 路径绑定到 grpc 服务方法上。
- PATCH 请求时，可以自动填充 FieldMask 字段。
- 自动生成 OpenAPI (Swagger) v2 API 定义。

[https://github.com/grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
