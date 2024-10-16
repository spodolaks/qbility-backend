package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	pb "github.com/spodolaks/qbility-backend/generated" // Correct import path for Protobuf-generated Go files
	"github.com/spodolaks/qbility-backend/server"       // Correct import path to the server package

	_ "github.com/go-sql-driver/mysql" // MySQL driver for Go
	"google.golang.org/grpc"
)

func main() {
	// Connect to the MySQL database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Listen on a TCP port for gRPC requests
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the Topic service on the gRPC server
	topicServer := server.NewTopicServer(db)
	pb.RegisterTopicServiceServer(grpcServer, topicServer)

	// Start the server
	fmt.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
