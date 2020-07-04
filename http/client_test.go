package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockClient(t *testing.T) {
	client := NewMockClient(
		WithMock("/route", &MockedResponseResult{
			Status: http.StatusConflict,
			Body:   "Conflict",
		}),
	)

	req, _ := http.NewRequest("GET", "/route", nil)

	res, err := client.Client().Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	assert.Contains(t, client.History(), "/route")
}

func TestMockClientAddMock(t *testing.T) {
	client := NewMockClient()

	client.AddMock("/my-route", &MockedResponseResult{
		Status: http.StatusBadRequest,
	})

	req, _ := http.NewRequest("GET", "/my-route", nil)

	res, err := client.Client().Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Contains(t, client.History(), "/my-route")
}

func TestMockClientDefaultTrip(t *testing.T) {
	client := NewMockClient()

	req, _ := http.NewRequest("GET", "/route", nil)

	res, err := client.Client().Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Contains(t, client.History(), "/route")
}
