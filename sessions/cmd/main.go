package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/VitoNaychev/elysium-challenge/crypto"
	"github.com/VitoNaychev/elysium-challenge/pgconfig"
	"github.com/VitoNaychev/elysium-challenge/rpc/sessions"
	"github.com/VitoNaychev/elysium-challenge/sessions/handler"
	"github.com/VitoNaychev/elysium-challenge/sessions/repository"
	"github.com/VitoNaychev/elysium-challenge/sessions/service"
	"google.golang.org/grpc"
)

func main() {
	pgConfig, err := pgconfig.InitFromEnv()
	if err != nil {
		log.Fatal("pgconfig InitFromEnv error: ", err)
	}

	userRepo, err := repository.NewPGUserRepository(context.Background(), pgConfig.GetConnectionString())
	if err != nil {
		log.Fatal("NewPGUserRepository error: ", err)
	}

	jwtConfig, err := crypto.InitJWTConfigFromEnv()
	if err != nil {
		log.Fatal("InitJWTConfigFromEnv error: ", err)
	}

	userService := service.NewUserService(jwtConfig, userRepo)

	userHTTPHandler := handler.NewUserHTTPHandler(userService)
	userRPCHandler := handler.NewUserRPCHandler(userService)

	httpServer := &http.Server{
		Addr:    "8080",
		Handler: userHTTPHandler,
	}

	rpcServer := grpc.NewServer()
	sessions.RegisterUserServer(rpcServer, userRPCHandler)

	go listenAndServeHTTP(httpServer, ":8080")
	go listenAndServeRPC(rpcServer, ":6060")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	sig := <-sigCh
	log.Printf("Received signal: %v. Shutting down...", sig)

	shutdownHTTPServer(httpServer)
	shutdownRPCServer(rpcServer)
}

func listenAndServeHTTP(server *http.Server, port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("net.Listen error: ", err)
	}

	log.Printf("Starting HTTP server on port %v...", port)
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func shutdownHTTPServer(server *http.Server) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
}

func listenAndServeRPC(server *grpc.Server, port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("net.Listen error: ", err)
	}

	go func() {
		log.Printf("Starting RPC server on port %v...", port)
		if err := server.Serve(listener); err != nil {
			log.Fatal("Serve error: ", err)
		}
	}()
}

func shutdownRPCServer(server *grpc.Server) {
	server.GracefulStop()
}
