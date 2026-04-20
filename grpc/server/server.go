package main

import (
	"context"
	"google.golang.org/grpc"
	"net"

	"go_test/grpc/proto"
)

type Server struct {
	proto.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hello " + req.Name}, nil
}
func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})

	lis, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	serveErr := g.Serve(lis)
	if serveErr != nil {
		panic("failed to start grpc:" + err.Error())
	}

}
