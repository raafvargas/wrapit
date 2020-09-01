package grpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	shutdown := make(chan os.Signal)
	server := NewServer(
		WithServerHost(":0"),
		WithServerShutdown(shutdown),
	)

	go server.Run(context.Background())

	time.Sleep(time.Second)

	_, err := NewClient(
		WithClientAddress(server.listener.Addr().String()),
	)

	assert.NoError(t, err)

	shutdown <- os.Interrupt
}
