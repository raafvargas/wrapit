package tracing

import (
	"strings"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

var (
	// AMQPPropagator ...
	AMQPPropagator propagation.HTTPPropagator = trace.B3{}
)

// AMQPSupplier ...
type AMQPSupplier amqp.Table

// Set ...
func (c AMQPSupplier) Set(key, val string) {
	key = strings.ToLower(key)

	c[key] = val
}

// Get ...
func (c AMQPSupplier) Get(key string) string {
	key = strings.ToLower(key)

	v := c[key]

	str, ok := v.(string)

	if ok {
		return str
	}

	return ""
}
