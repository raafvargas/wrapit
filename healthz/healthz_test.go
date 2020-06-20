package healthz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/raafvargas/wrapit/configuration"
	"github.com/raafvargas/wrapit/healthz"
	"github.com/raafvargas/wrapit/mongodb"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestHealthz(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	ctx.Request = httptest.NewRequest("GET", "/healthz", nil)

	healthz.HTTPHealthz()(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHealthzWithMongo(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	ctx.Request = httptest.NewRequest("GET", "/healthz", nil)

	cfg := new(configuration.Config)
	err := configuration.FromYAML("../tests/config.yaml", cfg)

	if err != nil {
		t.Fatal(err)
	}

	client, err := mongodb.Connect(context.Background(), uuid.New().String(), cfg.Mongo)

	if err != nil {
		t.Fatal(err)
	}

	handler := healthz.HTTPHealthz(
		healthz.WithMongo(client),
	)

	handler(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHealthzInvalidMongo(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	ctx.Request = httptest.NewRequest("GET", "/healthz", nil)

	client := &mongo.Client{}

	handler := healthz.HTTPHealthz(
		healthz.WithMongo(client),
	)

	handler(ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
