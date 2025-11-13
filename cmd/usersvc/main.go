package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	userv1 "usersvc/api/gen/go/user/v1"
	"usersvc/internal/adapters/driven/pub/stdout"
	"usersvc/internal/adapters/driven/repo/memory"
	"usersvc/internal/adapters/driving/grpc"
	appuser "usersvc/internal/app/user"
)

func main() {
	port := getenv("PORT", "8080")
	addr := ":" + port

	// Wire hex: Driving adapter (gRPC) -> UseCase -> Driven adapters (Repo, Publisher)
	repo := memory.New()
	pub := stdout.New()
	createUC := appuser.NewCreateUserUseCase(repo, pub)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}

	gs := grpc.NewServer()
	userv1.RegisterUserServiceServer(gs, grpcserver.New(createUC))

	hs := health.NewServer()
	healthpb.RegisterHealthServer(gs, hs)

	go func() {
		log.Printf("gRPC listening on %s", addr)
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Printf("shutting down...")
	gs.GracefulStop()
	log.Printf("bye")
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}
