package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	gintrace "go.opentelemetry.io/contrib/instrumentation/gin-gonic/gin"
)

// Config ...
type Config struct {
	Host string `yaml:"host"`
}

var (
	// DefaultHealthz ...
	DefaultHealthz = func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
)

// Server ...
type Server struct {
	Engine *gin.Engine

	Handlers    []gin.HandlerFunc
	NoRoute     []gin.HandlerFunc
	Metrics     http.Handler
	ServiceName string
	Host        string
	Controllers []Controller
	Healthz     gin.HandlerFunc
	Shutdown    chan os.Signal
}

// New ....
func New(opts ...Option) *Server {
	server := &Server{}
	server.Handlers = []gin.HandlerFunc{}
	server.NoRoute = []gin.HandlerFunc{}
	server.Controllers = []Controller{}
	server.Shutdown = make(chan os.Signal)
	server.Healthz = DefaultHealthz
	server.Metrics = promhttp.Handler()

	for _, opt := range opts {
		opt(server)
	}

	server.Engine = gin.Default()

	server.Engine.Use(
		gintrace.Middleware(server.ServiceName),
	)

	server.Engine.GET("healthz", server.Healthz)
	server.Engine.GET("metrics", func(ctx *gin.Context) {
		server.Metrics.ServeHTTP(ctx.Writer, ctx.Request)
	})

	server.Engine.NoRoute(server.NoRoute...)
	server.Engine.Use(server.Engine.Handlers...)

	for _, ctrl := range server.Controllers {
		ctrl.RegisterRoutes(&server.Engine.RouterGroup)
	}

	return server
}

// Run ...
func (server *Server) Run() error {
	srv := http.Server{
		Addr:    server.Host,
		Handler: server.Engine,
	}

	signal.Notify(server.Shutdown, os.Interrupt)

	go func() {
		<-server.Shutdown

		logrus.Info("waiting 5 seconds to stop the server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logrus.WithError(err).
				Error("shotdown error")
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
