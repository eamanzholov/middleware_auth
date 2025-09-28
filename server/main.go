package main

import (
	"context"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "middleware/proto/github.com/eamanzholov/middleware_auth"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

// ✅ Аутентификация через Bearer-токен
func authFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "no token provided")
	}
	if token != "super-secret" {
		return nil, status.Error(codes.PermissionDenied, "invalid token")
	}
	return ctx, nil
}

// ✅ Реализация метода Login
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Валидация автоматически вызовется через grpc_validator
	return &pb.LoginResponse{
		Token: fmt.Sprintf("fake-token-%d", time.Now().Unix()),
	}, nil
}

func main() {
	// Логгер
	logger, _ := zap.NewProduction()

	// gRPC сервер с middleware
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_zap.UnaryServerInterceptor(logger),   // Логирование
			grpc_auth.UnaryServerInterceptor(authFunc), // Аутентификация
			grpc_validator.UnaryServerInterceptor(),    // Валидация
		),
	)

	pb.RegisterAuthServiceServer(s, &server{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	logger.Info("🚀 gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}