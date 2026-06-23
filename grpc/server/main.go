package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"

	"grpc_sample/gen"
)

// proto の代わりに JSON でシリアライズするカスタムコーデック。
// "proto" という名前で登録することで gRPC のデフォルトコーデックを置き換える。
func init() {
	encoding.RegisterCodec(jsonCodec{})
}

type jsonCodec struct{}

func (jsonCodec) Name() string                       { return "proto" }
func (jsonCodec) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (jsonCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

// ─── UserService 実装 ─────────────────────────────────────────

type userServer struct {
	gen.UnimplementedUserServiceServer
	mu      sync.RWMutex
	users   map[string]*gen.User
	counter atomic.Int64
}

func newUserServer() *userServer {
	return &userServer{users: make(map[string]*gen.User)}
}

// Unary RPC: ユーザー作成
func (s *userServer) CreateUser(_ context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	id := fmt.Sprintf("user-%d", s.counter.Add(1))
	u := &gen.User{Id: id, Name: req.Name, Email: req.Email}

	s.mu.Lock()
	s.users[id] = u
	s.mu.Unlock()

	log.Printf("[CreateUser] id=%s name=%s email=%s", id, req.Name, req.Email)
	return &gen.CreateUserResponse{User: u}, nil
}

// Unary RPC: ユーザー取得
func (s *userServer) GetUser(_ context.Context, req *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	s.mu.RLock()
	u, ok := s.users[req.Id]
	s.mu.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "user %s not found", req.Id)
	}

	log.Printf("[GetUser] id=%s", req.Id)
	return &gen.GetUserResponse{User: u}, nil
}

// Server streaming RPC: ユーザー一覧をストリームで返す
func (s *userServer) ListUsers(_ *gen.ListUsersRequest, stream gen.UserService_ListUsersServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if err := stream.Send(&gen.ListUsersResponse{User: u}); err != nil {
			return err
		}
	}

	log.Printf("[ListUsers] sent %d users", len(s.users))
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer()
	gen.RegisterUserServiceServer(s, newUserServer())

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
