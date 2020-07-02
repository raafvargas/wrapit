package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/wrapit/auth"
	"github.com/stretchr/testify/assert"
)

func TestMockHandler(t *testing.T) {
	handler := auth.NewMock(true, map[string]interface{}{
		"a": "b",
	})

	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	handler.HTTP()(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "b", ctx.GetString("a"))
}

func TestMockHandlerUnauthorized(t *testing.T) {
	handler := auth.NewMock(false, nil)

	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	handler.HTTP()(ctx)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
