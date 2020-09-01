package grpc

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/trace"
)

type GRPCServerOption func(*GRPCServer)

type GRPClientOption func(*GRPCClient)

func WithClientAddress(address string) GRPClientOption {
	return func(s *GRPCClient) {
		s.address = address
	}
}

func WithClientTracer(tracer trace.Tracer) GRPClientOption {
	return func(s *GRPCClient) {
		s.tracer = tracer
	}
}

func WithServerHost(host string) GRPCServerOption {
	return func(s *GRPCServer) {
		s.host = host
	}
}

func WithServerLogger(logger *logrus.Logger) GRPCServerOption {
	return func(s *GRPCServer) {
		s.logger = logger
	}
}

func WithServerTracer(tracer trace.Tracer) GRPCServerOption {
	return func(s *GRPCServer) {
		s.tracer = tracer
	}
}

func WithServerShutdown(shutdown chan os.Signal) GRPCServerOption {
	return func(s *GRPCServer) {
		s.shutdown = shutdown
	}
}

func WithServerService(service GRPCService) GRPCServerOption {
	return func(s *GRPCServer) {
		s.services = append(s.services, service)
	}
}
