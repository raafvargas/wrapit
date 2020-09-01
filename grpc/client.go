package grpc

import (
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/instrumentation/grpctrace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GRPCClient struct {
	address string
	tracer  trace.Tracer
}

func NewClient(opts ...GRPClientOption) (*grpc.ClientConn, error) {
	config := &GRPCClient{
		tracer: global.Tracer("grpc"),
	}

	for _, o := range opts {
		o(config)
	}

	conn, err := grpc.Dial(
		config.address,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second * 30,
				Timeout:             time.Second * 10,
				PermitWithoutStream: true,
			},
		),

		grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(),
			grpctrace.UnaryClientInterceptor(config.tracer),
		),

		grpc.WithChainStreamInterceptor(
			grpc_retry.StreamClientInterceptor(),
			grpctrace.StreamClientInterceptor(config.tracer),
		),
	)
	return conn, err
}
