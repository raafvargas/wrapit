package rabbitmq_test

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/raafvargas/wrapit/configuration"
	"github.com/raafvargas/wrapit/rabbitmq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConsumerTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	config     *configuration.Config
	connection *rabbitmq.RabbitConnection

	queueName    string
	exchangeName string
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (s *ConsumerTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.config = new(configuration.Config)

	err := configuration.FromYAML("../tests/config.yaml", s.config)

	if err != nil {
		s.FailNow(err.Error())
	}

	connection, err := rabbitmq.NewConnection(s.config.RabbitMQ)

	if err != nil {
		s.FailNow(err.Error())
	}

	s.queueName = uuid.New().String()
	s.exchangeName = uuid.New().String()
	s.connection = connection

	s.connection.EnsureExchange(context.Background(), s.exchangeName)
	s.connection.EnsureQueue(context.Background(), s.queueName, s.exchangeName)
}

func (s *ConsumerTestSuite) TearDown() {
	s.connection.Close()
}

func (s *ConsumerTestSuite) TestConsumer() {
	message := struct {
		A string `json:"a"`
	}{
		A: "B",
	}

	called := make(chan bool, 1)
	done := make(chan interface{}, 1)

	consumer, err := rabbitmq.NewConsumer(
		s.connection,
		rabbitmq.WithMessageType(reflect.TypeOf(message)),
		rabbitmq.WithHandler(
			func(_ context.Context, message interface{}) error {
				m, ok := message.(*struct {
					A string `json:"a"`
				})

				s.assert.True(ok)
				s.assert.Equal("B", m.A)

				called <- true
				return nil
			},
		),
	)
	s.assert.NoError(err)

	go func() {
		if err := consumer.Consume(context.Background(), s.queueName); err != nil {
			s.T().Log(err)
			s.FailNow(err.Error())
		}

		done <- true
	}()

	conn, _ := rabbitmq.NewConnection(s.config.RabbitMQ)
	producer := rabbitmq.NewProducer(conn)
	defer conn.Close()

	producer.Publish(context.Background(), s.exchangeName, message)
	s.assert.NoError(err)

	s.assert.True(<-called)

	consumer.Shutdown <- os.Interrupt

	<-done
}

func (s *ConsumerTestSuite) TestConsumerInvalidMessage() {
	errCh := make(chan error, 1)

	consumer, err := rabbitmq.NewConsumer(
		s.connection,
		rabbitmq.WithMessageType(reflect.TypeOf(map[string]interface{}{})),
		rabbitmq.WithHandler(
			func(context.Context, interface{}) error { return nil }),
		rabbitmq.WithOnError(func(_ context.Context, err error) {
			errCh <- err
		}),
	)
	s.assert.NoError(err)

	go consumer.Consume(context.Background(), s.queueName)

	conn, _ := rabbitmq.NewConnection(s.config.RabbitMQ)
	producer := rabbitmq.NewProducer(conn)
	defer conn.Close()

	err = producer.Publish(context.Background(), s.exchangeName, "invalid_json")
	s.assert.NoError(err)

	s.assert.Error(<-errCh)
	s.assert.EqualError(<-errCh, rabbitmq.ErrInvalidMessageBody.Error())

	consumer.Shutdown <- os.Interrupt
}

func (s *ConsumerTestSuite) TestConsumerHandlerError() {
	errCh := make(chan error, 1)
	handlerErr := errors.New("handler error")

	message := struct {
		A string `json:"a"`
	}{
		A: "B",
	}

	consumer, err := rabbitmq.NewConsumer(
		s.connection,
		rabbitmq.WithMessageType(reflect.TypeOf(message)),
		rabbitmq.WithHandler(
			func(context.Context, interface{}) error {
				return handlerErr
			},
		),
		rabbitmq.WithOnError(func(_ context.Context, err error) {
			errCh <- err
		}),
	)
	s.assert.NoError(err)

	go consumer.Consume(context.Background(), s.queueName)

	conn, _ := rabbitmq.NewConnection(s.config.RabbitMQ)
	producer := rabbitmq.NewProducer(conn)
	defer conn.Close()

	err = producer.Publish(context.Background(), s.exchangeName, message)
	s.assert.NoError(err)

	s.assert.Error(<-errCh)
	s.assert.EqualError(<-errCh, handlerErr.Error())

	consumer.Shutdown <- os.Interrupt
}

func (s *ConsumerTestSuite) TestNewConsumerWithoutHandler() {
	_, err := rabbitmq.NewConsumer(s.connection)
	s.assert.Error(err)
	s.assert.EqualError(err, "handler must not be nil")
}

func (s *ConsumerTestSuite) TestNewConsumerWithoutMessageType() {
	_, err := rabbitmq.NewConsumer(s.connection, rabbitmq.WithHandler(func(context.Context, interface{}) error {
		return nil
	}))

	s.assert.Error(err)
	s.assert.EqualError(err, "messageType must not be nil")
}
