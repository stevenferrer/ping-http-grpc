package pingserver

import (
	"context"

	"github.com/stevenferrer/cmux-http-grpc/pb"
)

type PingServer struct {
	pb.UnimplementedPingServer
}

var _ pb.PingServer = (*PingServer)(nil)

func New() *PingServer {
	return &PingServer{}
}

func (s *PingServer) Ping(
	ctx context.Context,
	request *pb.PingRequest,
) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}
