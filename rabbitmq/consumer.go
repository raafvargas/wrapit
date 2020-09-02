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

// AMQPConsumer ...
type AMQPConsumer interface {
	Consume(ctx context.Context) error
}

// Consumer ...
type Consumer struct {
	Shutdown chan os.Signal
	stopped  chan error

	connection *RabbitConnection

	tracer trace.Tracer
	logger *logrus.Logger

	Queue        string
	Exchange     string
	MessageType  reflect.Type
	Prefetch     int
	Asynchronous int64
	Handler      AMQPHandler
	OnError      func(context.Context, error)
}

// NewConsumer ...
func NewConsumer(connection *RabbitConnection, options ...ConsumerOption) (*Consumer, error) {
	consumer := &Consumer{
		connection:   connection,
		logger:       logrus.New(),
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

	if consumer.Queue == "" || consumer.Exchange == "" {
		return nil, errors.New("queue and exchangee must not be empty")
	}

	return consumer, nil
}

// Consume ...
func (c *Consumer) Consume(ctx context.Context) error {
	if err := c.ensureQueue(ctx); err != nil {
		return err
	}

	signal.Notify(c.Shutdown, os.Interrupt)

	if err := c.connection.Channel.Qos(c.Prefetch, 0, false); err != nil {
		return err
	}

	c.stopped = make(chan error, 1)

	go c.createConsumer(ctx, c.Queue)

	return <-c.stopped
}

func (c *Consumer) createConsumer(ctx context.Context, queue string) {
	delivery, err := c.connection.Channel.Consume(queue, "", false, false, false, false, nil)

	if err != nil {
		c.stopped <- err
		return
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
	c.logger.WithField("queue", c.Queue).WithField("exchange", delivery.Exchange).
		Infof("start consuming message %s", delivery.MessageId)

	defer func() {
		if err := recover(); err != nil {
			c.logger.WithField("err", err).Errorf("consumer panicked")
			delivery.Reject(false)
		}
	}()

	ctx := tracing.AMQPPropagator.Extract(context.Background(), tracing.AMQPSupplier(delivery.Headers))

	ctx, span := c.tracer.Start(ctx, ConsumerOperationName,
		trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	message := reflect.New(c.MessageType).Interface()

	if err := json.Unmarshal(delivery.Body, message); err != nil {
		span.RecordError(ctx, err)
		c.logger.WithField("type", c.MessageType.String()).
			WithField("body", string(delivery.Body)).
			Warn("coldn't unmarshal message body")

		if err := delivery.Reject(false); err != nil {
			span.RecordError(ctx, err)
			c.logger.WithError(err).Error("nack error")
		}

		if c.OnError != nil {
			c.OnError(ctx, ErrInvalidMessageBody)
		}

		return
	}

	if err := c.Handler.Handle(ctx, message); err != nil {
		span.RecordError(ctx, err)

		c.logger.WithError(err).
			Error("consumer handler error")

		if err := delivery.Reject(false); err != nil {
			span.RecordError(ctx, err)
			c.logger.WithError(err).Error("nack error")
		}

		if c.OnError != nil {
			c.OnError(ctx, err)
		}

		return
	}

	if err := delivery.Ack(false); err != nil {
		span.RecordError(ctx, err)

		c.logger.WithError(err).
			Error("ack error")
		return
	}

	c.logger.Infof("finished message %s", delivery.MessageId)
}

func (c *Consumer) ensureQueue(ctx context.Context) error {
	if err := c.connection.EnsureExchange(ctx, c.Exchange); err != nil {
		return err
	}

	if err := c.connection.EnsureQueue(ctx, c.Queue, c.Exchange); err != nil {
		return err
	}

	return nil
}
