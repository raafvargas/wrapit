package grpc

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/api/trace"
)

func TestWithServerHost(t *testing.T) {
	grpcServer := &GRPCServer{}

	WithServerHost(":80")(grpcServer)

	assert.Equal(t, ":80", grpcServer.host)
}

func TestWithServerLogger(t *testing.T) {
	grpcServer := &GRPCServer{}

	WithServerLogger(logrus.New())(grpcServer)

	assert.NotNil(t, grpcServer.logger)
}

func TestWithServerTracer(t *testing.T) {
	grpcServer := &GRPCServer{}

	WithServerTracer(trace.NoopTracer{})(grpcServer)

	assert.NotNil(t, grpcServer.tracer)
}

func TestWithServerShutdown(t *testing.T) {
	grpcServer := &GRPCServer{}

	WithServerShutdown(make(chan os.Signal))(grpcServer)

	assert.NotNil(t, grpcServer.shutdown)
}

func TestWithServerService(t *testing.T) {
	grpcServer := &GRPCServer{}

	WithServerService(GRPCService(nil))(grpcServer)

	assert.Len(t, grpcServer.services, 1)
}

func TestWithClientAddress(t *testing.T) {
	grpcClient := &GRPCClient{}

	WithClientAddress(":9999")(grpcClient)

	assert.Equal(t, ":9999", grpcClient.address)
}

func TestWithClientTracer(t *testing.T) {
	grpcClient := &GRPCClient{}

	WithClientTracer(trace.NoopTracer{})(grpcClient)

	assert.NotNil(t, grpcClient.tracer)
}
