package rabbitmq_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raafvargas/wrapit/configuration"
	"github.com/raafvargas/wrapit/rabbitmq"
	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	cfg := new(configuration.Config)

	if err := configuration.FromYAML("../tests/config.yaml", cfg); err != nil {
		t.Fatal(err)
	}

	conn, err := rabbitmq.NewConnection(cfg.RabbitMQ)

	assert.NoError(t, err)
	assert.NoError(t, conn.Close())
}

func TestConnectionReconnect(t *testing.T) {
	cfg := new(configuration.Config)

	if err := configuration.FromYAML("../tests/config.yaml", cfg); err != nil {
		t.Fatal(err)
	}

	conn, err := rabbitmq.NewConnection(cfg.RabbitMQ)

	assert.NoError(t, err)
	assert.Contains(t, conn.Connection.Properties, rabbitmq.ConnectionIdentifierProperty)

	connID := conn.Connection.Properties[rabbitmq.ConnectionIdentifierProperty]

	err = conn.Channel.Close()
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	assert.False(t, conn.Connection.IsClosed())
	assert.NotEqual(t, connID, conn.Connection.Properties[rabbitmq.ConnectionIdentifierProperty])
	assert.NoError(t, conn.Close())
}

func TestEnsureQueue(t *testing.T) {
	cfg := new(configuration.Config)

	if err := configuration.FromYAML("../tests/config.yaml", cfg); err != nil {
		t.Fatal(err)
	}

	conn, err := rabbitmq.NewConnection(cfg.RabbitMQ)
	defer conn.Close()

	assert.NoError(t, err)

	queueName := uuid.New().String()
	exchangeName := uuid.New().String()

	err = conn.EnsureExchange(context.Background(), exchangeName)
	assert.NoError(t, err)

	err = conn.EnsureQueue(context.Background(), queueName, exchangeName)
	assert.NoError(t, err)
}

func TestEnsureQueueFail(t *testing.T) {
	cfg := new(configuration.Config)

	if err := configuration.FromYAML("../tests/config.yaml", cfg); err != nil {
		t.Fatal(err)
	}

	conn, err := rabbitmq.NewConnection(cfg.RabbitMQ)
	defer conn.Close()

	assert.NoError(t, err)

	queueName := uuid.New().String()
	exchangeName := uuid.New().String()

	err = conn.EnsureQueue(context.Background(), queueName, exchangeName)
	assert.Error(t, err)
}
