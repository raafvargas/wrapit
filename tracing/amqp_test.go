package tracing_test

import (
	"testing"

	"github.com/raafvargas/wrapit/tracing"
	"github.com/stretchr/testify/assert"
)

func TestAMQPSupplier(t *testing.T) {
	supplier := make(tracing.AMQPSupplier)

	supplier.Set("key", "value")

	assert.Equal(t, "value", supplier.Get("key"))

	supplier["key"] = 1

	assert.Equal(t, "", supplier.Get("key"))
}
