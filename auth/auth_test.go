package auth_test

import (
	"testing"

	"github.com/raafvargas/wrapit/auth"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	config := &auth.Config{
		Mock: auth.MockConfig{
			Enabled: true,
		},
	}

	handler := auth.NewHandler(config)
	assert.IsType(t, &auth.MockHandler{}, handler)

	config.Mock.Enabled = false
	handler = auth.NewHandler(config)
	assert.IsType(t, &auth.Auth0Handler{}, handler)
}
