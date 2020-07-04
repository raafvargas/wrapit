package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// AMQPMockProducer ...
type AMQPMockProducer struct {
	mock.Mock
}

// Publish ...
func (m *AMQPMockProducer) Publish(ctx context.Context, exchange string, message interface{}) error {
	args := m.Called(ctx, exchange, message)
	return args.Error(0)
}
