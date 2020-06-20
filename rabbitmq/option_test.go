package rabbitmq_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/raafvargas/wrapit/rabbitmq"

	"github.com/stretchr/testify/assert"
)

func TestWithPrefetch(t *testing.T) {
	consumer := &rabbitmq.Consumer{}

	rabbitmq.WithPrefetch(10)(consumer)

	assert.Equal(t, 10, consumer.Prefetch)
}

func TestWithAsynchronous(t *testing.T) {
	consumer := &rabbitmq.Consumer{}

	rabbitmq.WithAsynchronous(10)(consumer)

	assert.Equal(t, int64(10), consumer.Asynchronous)
}

func TestWithMessageType(t *testing.T) {
	consumer := &rabbitmq.Consumer{}

	rabbitmq.WithMessageType(reflect.TypeOf(0))(consumer)
	assert.Equal(t, "int", consumer.MessageType.String())
}

func TestWithHandler(t *testing.T) {
	consumer := &rabbitmq.Consumer{}

	called := make(chan bool, 1)

	rabbitmq.WithHandler(func(context.Context, interface{}) error {
		called <- true
		return nil
	})(consumer)

	err := consumer.Handler(context.Background(), nil)
	assert.NoError(t, err)
	assert.True(t, <-called)
}

func TestWithOnError(t *testing.T) {
	consumer := &rabbitmq.Consumer{}

	called := make(chan bool, 1)

	rabbitmq.WithOnError(func(context.Context, error) {
		called <- true
	})(consumer)

	consumer.OnError(context.Background(), nil)
	assert.True(t, <-called)
}
