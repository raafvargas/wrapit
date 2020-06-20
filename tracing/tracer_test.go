package tracing_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/raafvargas/wrapit/tracing"
	"github.com/stretchr/testify/assert"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

func TestRegisterJaeger(t *testing.T) {
	serviceName := uuid.New().String()
	spanName := uuid.New().String()
	eventName := uuid.New().String()

	called := make(chan bool, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		body, _ := ioutil.ReadAll(r.Body)
		assert.Contains(t, string(body), spanName)
		assert.Contains(t, string(body), eventName)
		assert.Contains(t, string(body), serviceName)

		called <- true
	}))
	defer ts.Close()

	provider, closer, err := tracing.Register(serviceName, &tracing.Config{
		JaegerURL: ts.URL,
	})
	assert.NoError(t, err)

	tracer := provider.Tracer(serviceName)
	ctx, span := tracer.Start(context.Background(), spanName)
	span.AddEvent(ctx, eventName)
	span.End()
	closer()

	assert.True(t, <-called)
}

func TestRegisterNoop(t *testing.T) {
	provider, _, err := tracing.Register(uuid.New().String(), &tracing.Config{
		JaegerDisabled: true,
	})
	assert.NoError(t, err)
	assert.IsType(t, &apitrace.NoopProvider{}, provider)
}
