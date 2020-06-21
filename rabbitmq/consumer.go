package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"reflect"

	"github.com/raafvargas/wrapit/tracing"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"golang.org/x/sync/semaphore"
)

var (
	// ConsumerOperationName ...
	ConsumerOperationName = "rabbitmq.consume"

	// ErrInvalidMessageBody ...
	ErrInvalidMessageBody = errors.New("couldn't unmarshal messagee")
)

// Consumer ...
type Consumer struct {
	Shutdown chan os.Signal
	stopped  chan error

	connection *RabbitConnection

	tracer trace.Tracer

	MessageType  reflect.Type
	Prefetch     int
	Asynchronous int64
	Handler      func(context.Context, interface{}) error
	OnError      func(context.Context, error)
}

// NewConsumer ...
func NewConsumer(connection *RabbitConnection, options ...ConsumerOption) (*Consumer, error) {
	consumer := &Consumer{
		connection:   connection,
		tracer:       global.Tracer(TracingTracerName),
		Asynchronous: 10,
		Prefetch:     100,
		Shutdown:     make(chan os.Signal, 1),
	}

	for _, o := range options {
		o(consumer)
	}

	if consumer.Handler == nil {
		return nil, errors.New("handler must not be nil")
	}

	if consumer.MessageType == nil {
		return nil, errors.New("messageType must not be nil")
	}

	return consumer, nil
}

// Consume ...
func (c *Consumer) Consume(ctx context.Context, queue string) error {
	signal.Notify(c.Shutdown, os.Interrupt)

	if err := c.connection.Channel.Qos(c.Prefetch, 0, false); err != nil {
		return err
	}

	c.stopped = make(chan error, 1)

	go c.createConsumer(ctx, queue)

	return <-c.stopped
}

func (c *Consumer) createConsumer(ctx context.Context, queue string) {
	delivery, err := c.connection.Channel.Consume(queue, "", false, false, false, false, nil)

	if err != nil {
		c.stopped <- err
	}

	sem := semaphore.NewWeighted(c.Asynchronous)

	for {
		select {
		case message := <-delivery:
			sem.Acquire(ctx, 1)
			go func() {
				defer sem.Release(1)
				c.handleDelivery(message)
			}()
		case sig := <-c.Shutdown:
			logrus.Infof("got sig %s. stopping consumers", sig.String())
			c.stopped <- nil
			return
		}
	}
}

func (c *Consumer) handleDelivery(delivery amqp.Delivery) {
	ctx, span := c.tracer.Start(context.Background(), ConsumerOperationName)
	defer span.End()

	logrus.WithField("headers", delivery.Headers).
		Info("message headers")

	ctx = tracing.AMQPPropagator.Extract(ctx, tracing.AMQPSupplier(delivery.Headers))

	message := reflect.New(c.MessageType).Interface()

	if err := json.Unmarshal(delivery.Body, message); err != nil {
		span.RecordError(ctx, err)
		logrus.WithField("type", c.MessageType.String()).
			WithField("body", string(delivery.Body)).
			Warn("coldn't unmarshal message body")

		if err := delivery.Nack(false, false); err != nil {
			span.RecordError(ctx, err)
			logrus.WithError(err).Error("nack error")
		}

		if c.OnError != nil {
			c.OnError(ctx, ErrInvalidMessageBody)
		}

		return
	}

	if err := c.Handler(ctx, message); err != nil {
		span.RecordError(ctx, err)

		logrus.WithError(err).
			Error("consumer handler error")

		if err := delivery.Nack(false, false); err != nil {
			span.RecordError(ctx, err)
			logrus.WithError(err).Error("nack error")
		}

		if c.OnError != nil {
			c.OnError(ctx, err)
		}

		return
	}

	if err := delivery.Ack(false); err != nil {
		span.RecordError(ctx, err)

		logrus.WithError(err).
			Error("ack error")
	}
}
