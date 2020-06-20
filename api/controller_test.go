package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/wrapit/api"
	"github.com/raafvargas/wrapit/contract"
	"github.com/stretchr/testify/assert"
)

func TestResolveError(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	err := contract.NewError(http.StatusBadRequest, "bad request")

	api.ResolveError(ctx, err)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestResolveErrorUnexpected(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	api.ResolveError(ctx, errors.New("unexpected"))

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestBindingError(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	api.BindingError(ctx, errors.New("binding"))

	assert.Equal(t, http.StatusBadRequest, res.Code)
}
