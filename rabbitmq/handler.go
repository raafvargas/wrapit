package rabbitmq

import "context"

// AMQPHandler ...
type AMQPHandler interface {
	Handle(context.Context, interface{}) error
}

// DefaultHandler ...
type DefaultHandler struct {
	handler func(context.Context, interface{}) error
}

// NewDefaultHandler ...
func NewDefaultHandler(handler func(context.Context, interface{}) error) *DefaultHandler {
	return &DefaultHandler{
		handler: handler,
	}
}

// Handle ...
func (h *DefaultHandler) Handle(ctx context.Context, payload interface{}) error {
	return h.handler(ctx, payload)
}
