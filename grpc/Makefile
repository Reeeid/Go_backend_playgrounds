setup:
	go mod tidy

server:
	go run ./server

client:
	go run ./client

# protoc を使って gen/ 以下を再生成したい場合:
#   brew install protobuf
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
generate:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/user.proto
