// app/app.go
package app

import (
	"fmt"
	"log"
	"net"
	"os"

	"final/db"
	"final/grpcservice"
	expressionpb "final/proto/expression"
	"final/routes"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Start() {
	db.InitDatabase()
	defer db.CloseDatabase()
	gormDB := db.GetDB()

	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = ":50051"
		}
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("gRPC listen on %s failed: %v", grpcPort, err)
		}

		grpcServer := grpc.NewServer()

		expressionpb.RegisterTaskServiceServer(
			grpcServer,
			grpcservice.NewServer(gormDB),
		)

		log.Printf("gRPC server listening on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC serve failed: %v", err)
		}
	}()

	r := gin.Default()
	routes.InitRoutes(r, gormDB)

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("HTTP server listening on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("HTTP serve failed: %v", err)
	}
}
