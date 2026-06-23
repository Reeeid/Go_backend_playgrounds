// Package gen はプロト定義から生成されるコードに相当する。
// 実際のプロジェクトでは protoc で自動生成する。
// 生成コマンド:
//   protoc --go_out=. --go-grpc_out=. proto/user.proto
package gen

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── メッセージ型 ──────────────────────────────────────────────

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest  struct{ Name, Email string }
type CreateUserResponse struct{ User *User }

type GetUserRequest  struct{ Id string }
type GetUserResponse struct{ User *User }

type ListUsersRequest  struct{}
type ListUsersResponse struct{ User *User }

// ─── サーバーインターフェース ─────────────────────────────────

type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	ListUsers(*ListUsersRequest, UserService_ListUsersServer) error
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer を埋め込むことで未実装メソッドがあっても
// コンパイルエラーにならず、呼ばれると Unimplemented を返す。
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "CreateUser not implemented")
}
func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetUser not implemented")
}
func (UnimplementedUserServiceServer) ListUsers(*ListUsersRequest, UserService_ListUsersServer) error {
	return status.Errorf(codes.Unimplemented, "ListUsers not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// ─── ストリーミングインターフェース ───────────────────────────

type UserService_ListUsersServer interface {
	Send(*ListUsersResponse) error
	grpc.ServerStream
}

type UserService_ListUsersClient interface {
	Recv() (*ListUsersResponse, error)
	grpc.ClientStream
}

type userServiceListUsersServer struct{ grpc.ServerStream }

func (x *userServiceListUsersServer) Send(m *ListUsersResponse) error {
	return x.ServerStream.SendMsg(m)
}

type userServiceListUsersClient struct{ grpc.ClientStream }

func (x *userServiceListUsersClient) Recv() (*ListUsersResponse, error) {
	m := new(ListUsersResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ─── クライアント ─────────────────────────────────────────────

type UserServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (UserService_ListUsersClient, error)
}

type userServiceClient struct{ cc grpc.ClientConnInterface }

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	if err := c.cc.Invoke(ctx, "/user.UserService/CreateUser", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	if err := c.cc.Invoke(ctx, "/user.UserService/GetUser", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (UserService_ListUsersClient, error) {
	stream, err := c.cc.NewStream(ctx, &_UserService_serviceDesc.Streams[0], "/user.UserService/ListUsers", opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceListUsersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// ─── サーバー登録 & ハンドラー ────────────────────────────────

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&_UserService_serviceDesc, srv)
}

func _UserService_CreateUser_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.UserService/CreateUser"}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUser_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.UserService/GetUser"}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ListUsers_Handler(srv any, stream grpc.ServerStream) error {
	m := new(ListUsersRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).ListUsers(m, &userServiceListUsersServer{stream})
}

var _UserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateUser", Handler: _UserService_CreateUser_Handler},
		{MethodName: "GetUser", Handler: _UserService_GetUser_Handler},
	},
	Streams: []grpc.StreamDesc{
		{StreamName: "ListUsers", Handler: _UserService_ListUsers_Handler, ServerStreams: true},
	},
	Metadata: "user.proto",
}
