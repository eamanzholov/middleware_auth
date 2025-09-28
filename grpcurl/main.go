package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/jhump/protoreflect/grpcreflect"
)

func main() {
	// подключаемся к gRPC серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	stub := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	client := grpcreflect.NewClient(ctx, stub)

	// список сервисов
	services, err := client.ListServices()
	if err != nil {
		log.Fatal(err)
	}

	for _, svc := range services {
		fmt.Println("🔹 Service:", svc)
		desc, _ := client.ResolveService(svc)

		// список методов
		for _, m := range desc.GetMethods() {
			fmt.Printf("   ▶ Method: %s (input: %s, output: %s)\n",
				m.GetName(),
				m.GetInputType().GetName(),
				m.GetOutputType().GetName(),
			)
		}
	}
}
