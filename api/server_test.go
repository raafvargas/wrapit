package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/wrapit/api"
	"github.com/raafvargas/wrapit/configuration"
	"github.com/stretchr/testify/assert"
)

type testController struct{}

func (*testController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("test", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
}

func TestServer(t *testing.T) {
	controller := new(testController)

	server := api.New(
		api.WithController(controller),
	)

	ts := httptest.NewServer(server.Engine)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/test", ts.URL))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestServerHealthz(t *testing.T) {
	server := api.New()

	ts := httptest.NewServer(server.Engine)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/healthz", ts.URL))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestServerMetrics(t *testing.T) {
	server := api.New()

	ts := httptest.NewServer(server.Engine)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/metrics", ts.URL))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestNoRoute(t *testing.T) {
	controller := new(testController)

	server := api.New(
		api.WithController(controller),
	)

	ts := httptest.NewServer(server.Engine)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/notfound", ts.URL))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestDefaultHealthz(t *testing.T) {
	res := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(res)

	api.DefaultHealthz(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHostBindingError(t *testing.T) {
	controller := new(testController)

	server := api.New(
		api.WithHost("invalid:host"),
		api.WithController(controller),
	)

	err := server.Run()
	assert.Error(t, err)
}

func TestGracefulShutdown(t *testing.T) {
	cfg := new(configuration.Config)

	err := configuration.FromYAML("../tests/config.yaml", cfg)
	assert.NoError(t, err)

	controller := new(testController)

	server := api.New(
		api.WithHost(cfg.API.Host),
		api.WithController(controller),
	)

	errCh := make(chan error, 1)

	go func() {
		errCh <- server.Run()
	}()

	server.Shutdown <- os.Interrupt

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("graceful shutdown failed")
	}
}
