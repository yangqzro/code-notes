package client

import (
	"context"
	"goexamples/gateway/openapi/internal/server"
	"goexamples/gateway/openapi/proto"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PostAutoFillFieldMask(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	bodyReader, ok := ctx.Value(server.BodyReaderContext).(io.Reader)
	if !ok || method != "/user.UserService/CreateUser" {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	if c, ok := req.(*proto.CreateUserRequest); ok && (c.CreateMask == nil || len(c.CreateMask.Paths) == 0) {
		fieldMask, err := runtime.FieldMaskFromRequestBody(bodyReader, c.User)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "%v", err)
		}
		c.CreateMask = fieldMask
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}
