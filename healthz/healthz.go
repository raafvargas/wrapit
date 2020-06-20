package healthz

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/wrapit/configuration"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Healthz ...
type Healthz struct {
	config *configuration.Config
	checks []func(context.Context, *Healthz) error
}

// HealtzOption ...
type HealtzOption func(*Healthz)

// WithMongo ...
func WithMongo(client *mongo.Client) HealtzOption {
	return func(h *Healthz) {
		h.checks = append(h.checks, func(ctx context.Context, healthz *Healthz) error {
			return client.Ping(ctx, readpref.Primary())
		})
	}
}

// NewHealthz ...
func NewHealthz(options ...HealtzOption) *Healthz {
	h := new(Healthz)
	h.checks = []func(context.Context, *Healthz) error{}

	for _, o := range options {
		o(h)
	}

	return h
}

// HTTPHealthz ...
func HTTPHealthz(options ...HealtzOption) gin.HandlerFunc {
	h := NewHealthz(options...)
	return h.HTTP()
}

// Healthz ...
func (h *Healthz) Healthz(ctx context.Context) error {
	wg := new(sync.WaitGroup)

	errCh := make(chan error, len(h.checks))
	doneCh := make(chan bool, len(h.checks))

	for _, check := range h.checks {
		wg.Add(1)
		go func(c func(context.Context, *Healthz) error) {
			defer wg.Done()
			if err := c(ctx, h); err != nil {
				errCh <- err
			}
		}(check)
	}

	go func() {
		wg.Wait()
		doneCh <- true
	}()

	<-doneCh

	close(errCh)
	close(doneCh)

	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}

// HTTP ...
func (h *Healthz) HTTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := h.Healthz(ctx.Request.Context()); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}
