package main

import (
	"context"
	"fmt"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "middleware/proto/github.com/eamanzholov/middleware_auth"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithMax(3),
				grpc_retry.WithPerRetryTimeout(2*time.Second),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Добавляем токен
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "bearer super-secret")

	resp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "admin",
		Password: "123456",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Token:", resp.Token)
}