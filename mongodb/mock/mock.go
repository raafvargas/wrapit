package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MongoRepository ...
type MongoRepository struct {
	mock.Mock
}

// Insert ...
func (m *MongoRepository) Insert(ctx context.Context, document interface{}) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

// Update ...
func (m *MongoRepository) Update(ctx context.Context, id interface{}, document interface{}) error {
	args := m.Called(ctx, id, document)
	return args.Error(0)
}

// FindByID ...
func (m *MongoRepository) FindByID(ctx context.Context, id interface{}) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}
