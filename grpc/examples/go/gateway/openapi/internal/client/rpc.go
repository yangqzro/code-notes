package client

import (
	"goexamples/gateway/openapi/proto"

	"google.golang.org/grpc"
)

type UserRPCClient struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
}

func (cli *UserRPCClient) RawClient() proto.UserServiceClient {
	return cli.client
}

func (cli *UserRPCClient) Close() error {
	return cli.conn.Close()
}

// func (cli *UserRPCClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserResponse, error) {
// }

// func (cli *UserRPCClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
// }

// func (cli *UserRPCClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error) {
// }

// func (cli *UserRPCClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
// }

// func (cli *UserRPCClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ListUsersResponse], error) {
// }

func NewUserRPCClientFromConn(conn *grpc.ClientConn) *UserRPCClient {
	return &UserRPCClient{conn: conn, client: proto.NewUserServiceClient(conn)}
}

func NewUserRPCClient(addr string, opts ...grpc.DialOption) *UserRPCClient {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		panic(err)
	}
	return NewUserRPCClientFromConn(conn)
}
