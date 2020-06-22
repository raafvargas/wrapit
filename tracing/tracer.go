package tracing

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"

	apitrace "go.opentelemetry.io/otel/api/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Config ...
type Config struct {
	JaegerURL      string `yaml:"jaeger_url"`
	JaegerDisabled bool   `yaml:"jaeger_disabled"`
}

// Register ...
func Register(serviceName string, config *Config) (apitrace.Provider, func(), error) {
	if config.JaegerDisabled {
		return NoopTracer()
	}

	tp, flush, err := jaeger.NewExportPipeline(
		jaeger.WithCollectorEndpoint(config.JaegerURL),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: serviceName,
		}),
		jaeger.RegisterAsGlobal(),
		jaeger.WithOnError(func(err error) {
			logrus.
				WithError(err).
				Errorf("error sending traces to jaeger")
		}),
		jaeger.WithBufferMaxCount(10),
		jaeger.WithSDK(
			&sdktrace.Config{
				DefaultSampler: sdktrace.AlwaysSample(),
			},
		),
	)

	global.SetTraceProvider(tp)

	return tp, flush, err
}

// NoopTracer ...
func NoopTracer() (apitrace.Provider, func(), error) {
	provider := &apitrace.NoopProvider{}

	return provider, func() {}, nil
}
