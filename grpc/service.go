package grpc

import (
	"google.golang.org/grpc"
)

type GRPCService interface {
	RegisterService(s *grpc.Server)
}
