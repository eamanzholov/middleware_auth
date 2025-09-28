package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
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
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_zap.UnaryServerInterceptor(logger),    // Логирование
			grpc_auth.UnaryServerInterceptor(authFunc), // Аутентификация
			grpc_validator.UnaryServerInterceptor(),    // Валидация
		),
	)

	// Регистрируем сервис
	pb.RegisterAuthServiceServer(grpcServer, &server{})

	// Включаем reflection (для grpcurl и других клиентов)
	reflection.Register(grpcServer)

	// Стартуем сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger.Info("🚀 gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}