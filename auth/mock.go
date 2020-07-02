package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MockHandler ...
type MockHandler struct {
	authenticated bool
	claims        map[string]interface{}
}

// NewMock ...
func NewMock(authenticated bool, claims map[string]interface{}) Handler {
	return &MockHandler{
		claims:        claims,
		authenticated: authenticated,
	}
}

// HTTP ...
func (h *MockHandler) HTTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !h.authenticated {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		for k, v := range h.claims {
			ctx.Set(k, v)
		}
	}
}
