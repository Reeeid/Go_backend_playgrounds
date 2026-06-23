package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"

	"grpc_sample/gen"
)

func init() {
	encoding.RegisterCodec(jsonCodec{})
}

type jsonCodec struct{}

func (jsonCodec) Name() string                       { return "proto" }
func (jsonCodec) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (jsonCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ── Unary: CreateUser ──────────────────────────────────────
	r1, err := client.CreateUser(ctx, &gen.CreateUserRequest{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		log.Fatalf("CreateUser: %v", err)
	}
	log.Printf("Created: %+v", r1.User)

	_, err = client.CreateUser(ctx, &gen.CreateUserRequest{Name: "Bob", Email: "bob@example.com"})
	if err != nil {
		log.Fatalf("CreateUser: %v", err)
	}

	// ── Unary: GetUser ─────────────────────────────────────────
	r2, err := client.GetUser(ctx, &gen.GetUserRequest{Id: r1.User.Id})
	if err != nil {
		log.Fatalf("GetUser: %v", err)
	}
	log.Printf("Got:     %+v", r2.User)

	// ── Server streaming: ListUsers ────────────────────────────
	stream, err := client.ListUsers(ctx, &gen.ListUsersRequest{})
	if err != nil {
		log.Fatalf("ListUsers: %v", err)
	}
	log.Println("All users:")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream: %v", err)
		}
		log.Printf("  %+v", resp.User)
	}
}
