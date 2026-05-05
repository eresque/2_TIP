package main

import (
	"log"
	"net"

	"example.com/pz2-grpc/gen/studentpb"
	"example.com/pz2-grpc/internal/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	repo := student.NewRepository()
	service := student.NewService(repo)

	server := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(server, service)
	reflection.Register(server)

	log.Println("gRPC server started on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
