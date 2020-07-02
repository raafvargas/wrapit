package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/wrapit/auth"

	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	ctx.Set("permissions", []interface{}{
		"read",
	})

	auth.Authorize("read")(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestAuthorizeWithoutPermission(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	ctx.Set("permissions", []interface{}{
		"read",
	})

	auth.Authorize("write")(ctx)

	assert.Equal(t, http.StatusForbidden, res.Code)
}

func TestAuthorizeWithoutPermissions(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	auth.Authorize("read")(ctx)

	assert.Equal(t, http.StatusForbidden, res.Code)
}
