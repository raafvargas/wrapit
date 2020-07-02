package auth

import (
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
)

// Auth0Handler ...
type Auth0Handler struct {
	validator *auth0.JWTValidator
}

// NewAuth0Handler ...
func NewAuth0Handler(config *Config) *Auth0Handler {
	handler := &Auth0Handler{}

	handler.validator = auth0.NewValidator(
		auth0.NewConfiguration(
			auth0.NewJWKClient(
				auth0.JWKClientOptions{
					URI: config.JWKS,
				},
				nil),
			config.Audience,
			config.Tenant,
			jose.RS256,
		),
		nil,
	)

	return handler
}

// HTTP ...
func (h *Auth0Handler) HTTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.validator.ValidateRequest(ctx.Request)

		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		claims := make(map[string]interface{})
		if err := h.validator.Claims(ctx.Request, token, &claims); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		for key, value := range claims {
			ctx.Set(key, value)
		}

		ctx.Next()
	}
}
