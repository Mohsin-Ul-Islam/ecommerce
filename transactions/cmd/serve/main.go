package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/mohsin-ul-islam/ecommerce/transactions/api/v1"
	pb "github.com/mohsin-ul-islam/ecommerce/transactions/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	conn, err := pgx.Connect(context.Background(), "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		log.Println("cannot connect to database")
		log.Fatal(err.Error())
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))
	if err != nil {
		log.Println("cannot listen on :8081")
		log.Fatal(err.Error())
	}

	server := grpc.NewServer()
	service := v1.NewTransactionService(conn)
	pb.RegisterTransactionServiceServer(server, service)
	reflection.Register(server)

	err = server.Serve(listener)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot serve on :%s", port))
	}

	log.Println(fmt.Sprintf("listening on :%s", port))
}
