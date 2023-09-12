package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/mohsin-ul-islam/ecommerce/customers/api/v1"
	pb "github.com/mohsin-ul-islam/ecommerce/customers/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		log.Println(fmt.Sprintf("cannot listen on :%s"), port)
		log.Fatal(err.Error())
	}

	server := grpc.NewServer()
	service := v1.NewCustomerService(conn)
	pb.RegisterCustomerServiceServer(server, service)
	reflection.Register(server)

	err = server.Serve(listener)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot serve on :%s", port))
	}
}
