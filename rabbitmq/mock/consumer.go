package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// AMQPMockConsumer ...
type AMQPMockConsumer struct {
	mock.Mock
}

// Consume ...
func (m *AMQPMockConsumer) Consume(ctx context.Context, queue string) error {
	args := m.Called(ctx, queue)
	return args.Error(0)
}
