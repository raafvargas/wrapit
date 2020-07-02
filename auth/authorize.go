package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authorize ...
func Authorize(scope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permissions, exists := ctx.Get("permissions")

		if !exists {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		for _, permission := range permissions.([]interface{}) {
			if permission == scope {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatus(http.StatusForbidden)
	}
}
