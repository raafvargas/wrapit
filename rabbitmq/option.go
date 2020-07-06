package rabbitmq

import (
	"context"
	"reflect"
)

// ConsumerOption ...
type ConsumerOption func(*Consumer)

// WithPrefetch ...
func WithPrefetch(prefetch int) ConsumerOption {
	return func(c *Consumer) {
		c.Prefetch = prefetch
	}
}

// WithAsynchronous ...
func WithAsynchronous(asynchronous int64) ConsumerOption {
	return func(c *Consumer) {
		c.Asynchronous = asynchronous
	}
}

// WithMessageType ...
func WithMessageType(t reflect.Type) ConsumerOption {
	return func(c *Consumer) {
		c.MessageType = t
	}
}

// WithHandler ...
func WithHandler(handler AMQPHandler) ConsumerOption {
	return func(c *Consumer) {
		c.Handler = handler
	}
}

// WithOnError ...
func WithOnError(onError func(context.Context, error)) ConsumerOption {
	return func(c *Consumer) {
		c.OnError = onError
	}
}

// WithQueue ...
func WithQueue(queue string) ConsumerOption {
	return func(c *Consumer) {
		c.Queue = queue
	}
}

// WithExchange ...
func WithExchange(exchange string) ConsumerOption {
	return func(c *Consumer) {
		c.Exchange = exchange
	}
}
