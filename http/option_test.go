package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMock(t *testing.T) {
	client := NewMockClient()

	WithMock("/route", &MockedResponseResult{
		Status: http.StatusOK,
	})(client)

	assert.Contains(t, client.result, "/route")
}

func TestWithRoundTrip(t *testing.T) {
	client := NewMockClient()

	WithRoundTrip(nil)(client)

	assert.Nil(t, client.roundTrip)
}
