package service

import (
	"context"
	gen "sona/gen"
	"sona/gen/sonav1connect"

	"connectrpc.com/connect"
)

type HelloServer struct {
	sonav1connect.UnimplementedHelloServiceHandler
}

func NewHelloServer() *HelloServer {
	return &HelloServer{}
}

func (s *HelloServer) Hello(
	ctx context.Context,
	req *connect.Request[gen.HelloRequest],
) (*connect.Response[gen.HelloResponse], error) {
	message := "Hello, " + req.Msg.Name + "!"
	return connect.NewResponse(&gen.HelloResponse{
		Message: message,
	}), nil
}
