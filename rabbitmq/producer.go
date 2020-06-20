package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/raafvargas/wrapit/tracing"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

var (
	// PublisherOperationName ...
	PublisherOperationName = "rabbitmq.publish"
)

// Producer ...
type Producer struct {
	connection  *RabbitConnection
	propagators propagation.Propagators
	tracer      trace.Tracer
}

// NewProducer ...
func NewProducer(connection *RabbitConnection) *Producer {
	return &Producer{
		connection:  connection,
		propagators: global.Propagators(),
		tracer:      global.Tracer(TracingTracerName),
	}
}

// Publish ...
func (p *Producer) Publish(ctx context.Context, exchange string, message interface{}) error {
	ctx, span := p.tracer.Start(ctx, PublisherOperationName)
	defer span.End()

	headers := make(tracing.AMQPSupplier)
	propagation.InjectHTTP(ctx, p.propagators, headers)

	data, err := json.Marshal(message)

	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	span.SetAttribute("message.body", string(data))

	return p.connection.Channel.Publish(exchange, "", true, false, amqp.Publishing{
		DeliveryMode: 2,
		Body:         data,
		ContentType:  "application/json",
		Headers:      amqp.Table(headers),
	})
}
