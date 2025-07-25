package server

import (
	"context"
	"fmt"
	"goexamples/gateway/openapi/internal/model"
	"goexamples/gateway/openapi/proto"
	"goexamples/utils"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type UserRPCServer struct {
	proto.UnimplementedUserServiceServer
	server *grpc.Server
	model  *model.UserModel
}

func (srv *UserRPCServer) SetModel(model *model.UserModel) {
	srv.model = model
}

func (srv *UserRPCServer) RawServer() *grpc.Server {
	return srv.server
}

func (srv *UserRPCServer) GetServiceInfo() map[string]grpc.ServiceInfo {
	return srv.server.GetServiceInfo()
}

func (srv *UserRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.server.ServeHTTP(w, r)
}

func (srv *UserRPCServer) Listen(port int, listener func(net.Listener)) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	if listener != nil {
		listener(lis)
	}
	return srv.server.Serve(lis)
}

func (srv *UserRPCServer) CreateUser(_ context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	log.Printf("server.CreateUser received: request=%s, user=%s\n", utils.String(req), utils.String(req.GetUser()))
	user := req.GetUser()
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no user")
	}
	user = srv.model.Create(user, req.GetCreateMask())
	return &proto.CreateUserResponse{User: user}, nil
}

func (srv *UserRPCServer) DeleteUser(_ context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	log.Printf("server.DeleteUser received: request=%s\n", utils.String(req))
	if user, ok := srv.model.Delete(req.GetId()); ok {
		return &proto.DeleteUserResponse{User: user}, nil
	}
	return nil, status.Errorf(codes.FailedPrecondition, "can not delete user")
}

func (srv *UserRPCServer) UpdateUser(_ context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	log.Printf("server.UpdateUser received: request=%s, user=%s\n", utils.String(req), utils.String(req.GetUser()))
	if user, ok := srv.model.Update(req.GetUser(), req.GetUpdateMask()); ok {
		return &proto.UpdateUserResponse{User: user}, nil
	}
	return nil, status.Errorf(codes.FailedPrecondition, "can not update user")
}

func (srv *UserRPCServer) GetUser(_ context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	log.Printf("server.GetUser received: request=%s\n", utils.String(req))
	if user, ok := srv.model.Get(req.GetId()); ok {
		return &proto.GetUserResponse{User: user}, nil
	}
	return nil, status.Errorf(codes.FailedPrecondition, "can not get user")
}

func (srv *UserRPCServer) ListUsers(_ *proto.ListUsersRequest, sout grpc.ServerStreamingServer[proto.ListUsersResponse]) error {
	log.Println("server.ListUsers received request")
	for _, user := range srv.model.List() {
		if err := sout.Send(&proto.ListUsersResponse{User: user}); err != nil {
			return status.Errorf(codes.Internal, "failed to send user: %v", err)
		} else {
			log.Printf("server.ListUsers send: user=%v\n", utils.String(user))
		}
	}
	return nil
}

func NewUserRPCServer(mode string, opts ...grpc.ServerOption) *UserRPCServer {
	server := grpc.NewServer(opts...)
	srv := &UserRPCServer{server: server}
	if mode == "dev" {
		reflection.Register(server)
	}
	proto.RegisterUserServiceServer(server, srv)
	return srv
}
