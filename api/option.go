package api

import (
	"github.com/gin-gonic/gin"
)

// Option wrapps all server configurations
type Option func(server *Server)

// WithHost ...
func WithHost(host string) Option {
	return func(server *Server) {
		server.Host = host
	}
}

// WithController ...
func WithController(controller Controller) Option {
	return func(server *Server) {
		server.Controllers = append(server.Controllers, controller)
	}
}

// WithServiceName ...
func WithServiceName(serviceName string) Option {
	return func(server *Server) {
		server.ServiceName = serviceName
	}
}

// WithHealthz ...
func WithHealthz(healthz gin.HandlerFunc) Option {
	return func(server *Server) {
		server.Healthz = healthz
	}
}

// WithHandler ...
func WithHandler(handler gin.HandlerFunc) Option {
	return func(server *Server) {
		server.Handlers = append(server.Handlers, handler)
	}
}

// WithNoRoute ...
func WithNoRoute(handler gin.HandlerFunc) Option {
	return func(server *Server) {
		server.NoRoute = append(server.NoRoute, handler)
	}
}
