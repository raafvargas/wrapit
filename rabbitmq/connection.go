package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	// TracingTracerName ...
	TracingTracerName = "rabbitmq"

	// DefaultReconnectDelay ...
	DefaultReconnectDelay = time.Second * 5

	// ConnectionIdentifierProperty ...
	ConnectionIdentifierProperty = "id"

	// DeadLetterSufix ..
	DeadLetterSufix = "dlq"

	// DeadLetterExchange ...
	DeadLetterExchange = "dead-letter"
)

// RabbitConfig ...
type RabbitConfig struct {
	URL string `yaml:"url"`
}

// RabbitConnection ...
type RabbitConnection struct {
	url                   string
	connected             bool
	reconnectSecondsDelay time.Duration
	mutex                 *sync.Mutex
	closed                chan *amqp.Error
	done                  chan interface{}

	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewConnection ...
func NewConnection(config *RabbitConfig) (*RabbitConnection, error) {
	rc := &RabbitConnection{
		url:                   config.URL,
		reconnectSecondsDelay: DefaultReconnectDelay,
		mutex:                 new(sync.Mutex),
		done:                  make(chan interface{}, 1),
	}

	if err := rc.connect(); err != nil {
		return nil, err
	}

	rc.connected = true

	go rc.handleReconnect()

	return rc, nil
}

// EnsureQueue ...
func (rc *RabbitConnection) EnsureQueue(ctx context.Context, queueName, exchangeName string) error {
	dlqQueue := fmt.Sprintf("%s.%s", queueName, DeadLetterSufix)

	attributes := make(amqp.Table)
	attributes["x-dead-letter-exchange"] = DeadLetterExchange
	attributes["x-dead-letter-routing-key"] = dlqQueue

	if err := rc.createQueue(ctx, queueName, "", exchangeName, attributes); err != nil {
		return err
	}

	_, err := rc.Channel.QueueDeclare(dlqQueue, true, false, false, false, nil)
	return err
}

// EnsureExchange ...
func (rc *RabbitConnection) EnsureExchange(ctx context.Context, exchangeName string) error {
	if err := rc.Channel.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil); err != nil {
		return err
	}

	return rc.Channel.ExchangeDeclare(DeadLetterExchange, "fanout", true, false, false, false, nil)
}

// Close ...
func (rc *RabbitConnection) Close() error {
	rc.done <- true

	if !rc.connected {
		return nil
	}

	rc.connected = false

	if err := rc.Channel.Close(); err != nil {
		return err
	}

	return rc.Connection.Close()
}

func (rc *RabbitConnection) connect() error {
	conn, err := amqp.Dial(rc.url)

	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	rc.Connection = conn
	rc.Channel = ch
	rc.closed = make(chan *amqp.Error)
	rc.Channel.NotifyClose(rc.closed)

	rc.Connection.Properties[ConnectionIdentifierProperty] = uuid.New().String()

	return nil
}

func (rc *RabbitConnection) handleReconnect() {
	for {
		rc.mutex.Lock()
		for !rc.connected {
			if err := rc.connect(); err != nil {
				logrus.WithError(err).
					Error("error reconnecting to rabbitmq")

				time.Sleep(rc.reconnectSecondsDelay)
				continue
			}
			rc.connected = true
		}
		rc.mutex.Unlock()

		select {
		case <-rc.done:
			return
		case err := <-rc.closed:
			rc.mutex.Lock()
			rc.connected = false
			rc.mutex.Unlock()
			logrus.WithError(err).
				Warnf("got an connection closed notification")
		}
	}
}

func (rc *RabbitConnection) createQueue(ctx context.Context, queueName, routingKey, exchangeName string, attributes amqp.Table) error {
	_, err := rc.Channel.QueueDeclare(queueName, true, false, false, false, attributes)

	if err != nil {
		return err
	}

	return rc.Channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
}
