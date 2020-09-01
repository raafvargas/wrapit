package grpc

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerRun(t *testing.T) {
	errCh := make(chan error)
	shutdown := make(chan os.Signal)

	grpcServer := NewServer(
		WithServerHost(":0"),
		WithServerShutdown(shutdown),
	)

	go func() {
		err := grpcServer.Run(context.Background())
		errCh <- err
	}()

	shutdown <- os.Interrupt

	assert.NoError(t, <-errCh)
}

func TestServerInvalidHost(t *testing.T) {
	grpcServer := NewServer(
		WithServerHost("invalidhost"),
	)

	err := grpcServer.Run(context.Background())

	assert.Error(t, err)
}
