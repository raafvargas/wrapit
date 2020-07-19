package api_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/wrapit/api"
	"github.com/stretchr/testify/assert"
)

func TestWithHost(t *testing.T) {
	server := new(api.Server)

	api.WithHost("host")(server)

	assert.Equal(t, "host", server.Host)
}

func TestWithController(t *testing.T) {
	controller := new(testController)
	server := new(api.Server)

	api.WithController(controller)(server)

	assert.Contains(t, server.Controllers, controller)
}

func TestWithServiceName(t *testing.T) {
	server := new(api.Server)

	api.WithServiceName("service")(server)

	assert.Contains(t, server.ServiceName, "service")
}

func TestWithHealthz(t *testing.T) {
	server := new(api.Server)

	called := make(chan bool, 1)
	api.WithHealthz(func(*gin.Context) {
		called <- true
	})(server)

	server.Healthz(nil)

	assert.True(t, <-called)
}

func TestWithHandler(t *testing.T) {
	server := new(api.Server)

	api.WithHandler(func(*gin.Context) {})(server)
	assert.Len(t, server.Handlers, 1)
}

func TestWithNoRoute(t *testing.T) {
	server := new(api.Server)

	api.WithNoRoute(func(*gin.Context) {})(server)
	assert.Len(t, server.NoRoute, 1)
}
