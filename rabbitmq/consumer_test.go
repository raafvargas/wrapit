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
		rabbitmq.WithQueue(s.queueName),
		rabbitmq.WithExchange(s.exchangeName),
		rabbitmq.WithMessageType(reflect.TypeOf(message)),
		rabbitmq.WithHandler(
			rabbitmq.NewDefaultHandler(
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
		),
	)
	s.assert.NoError(err)

	go func() {
		if err := consumer.Consume(context.Background()); err != nil {
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

func (s *ConsumerTestSuite) TestConsumerPanic() {
	message := struct {
		A string `json:"a"`
	}{
		A: "B",
	}

	called := make(chan bool, 1)
	done := make(chan interface{}, 1)

	consumer, err := rabbitmq.NewConsumer(
		s.connection,
		rabbitmq.WithQueue(s.queueName),
		rabbitmq.WithExchange(s.exchangeName),
		rabbitmq.WithMessageType(reflect.TypeOf(message)),
		rabbitmq.WithHandler(
			rabbitmq.NewDefaultHandler(
				func(_ context.Context, message interface{}) error {
					called <- true
					panic("got some error")
				},
			),
		),
	)
	s.assert.NoError(err)

	go func() {
		if err := consumer.Consume(context.Background()); err != nil {
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
		rabbitmq.WithQueue(s.queueName),
		rabbitmq.WithExchange(s.exchangeName),
		rabbitmq.WithMessageType(reflect.TypeOf(map[string]interface{}{})),
		rabbitmq.WithHandler(
			rabbitmq.NewDefaultHandler(
				func(context.Context, interface{}) error {
					return nil
				},
			),
		),
		rabbitmq.WithOnError(func(_ context.Context, err error) {
			errCh <- err
		}),
	)
	s.assert.NoError(err)

	go consumer.Consume(context.Background())

	conn, _ := rabbitmq.NewConnection(s.config.RabbitMQ)
	producer := rabbitmq.NewProducer(conn)
	defer conn.Close()

	err = producer.Publish(context.Background(), s.exchangeName, "invalid_json")
	s.assert.NoError(err)

	consumerErr := <-errCh

	s.assert.Error(consumerErr)
	s.assert.EqualError(consumerErr, rabbitmq.ErrInvalidMessageBody.Error())

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
		rabbitmq.WithQueue(s.queueName),
		rabbitmq.WithExchange(s.exchangeName),
		rabbitmq.WithMessageType(reflect.TypeOf(message)),
		rabbitmq.WithHandler(
			rabbitmq.NewDefaultHandler(
				func(context.Context, interface{}) error {
					return handlerErr
				},
			),
		),
		rabbitmq.WithOnError(func(_ context.Context, err error) {
			errCh <- err
		}),
	)
	s.assert.NoError(err)

	go consumer.Consume(context.Background())

	conn, _ := rabbitmq.NewConnection(s.config.RabbitMQ)
	producer := rabbitmq.NewProducer(conn)
	defer conn.Close()

	err = producer.Publish(context.Background(), s.exchangeName, message)
	s.assert.NoError(err)

	consumerErr := <-errCh

	s.assert.Error(consumerErr)
	s.assert.EqualError(consumerErr, handlerErr.Error())

	consumer.Shutdown <- os.Interrupt
}

func (s *ConsumerTestSuite) TestNewConsumerWithoutHandler() {
	_, err := rabbitmq.NewConsumer(s.connection)
	s.assert.Error(err)
	s.assert.EqualError(err, "handler must not be nil")
}

func (s *ConsumerTestSuite) TestNewConsumerWithoutMessageType() {
	_, err := rabbitmq.NewConsumer(
		s.connection,
		rabbitmq.WithHandler(
			rabbitmq.NewDefaultHandler(
				func(context.Context, interface{}) error {
					return nil
				},
			),
		),
	)

	s.assert.Error(err)
	s.assert.EqualError(err, "messageType must not be nil")
}

func (s *ConsumerTestSuite) TestNewConsumerWithoutExchangeAndQueue() {
	_, err := rabbitmq.NewConsumer(
		s.connection,
		rabbitmq.WithMessageType(reflect.TypeOf("")),
		rabbitmq.WithHandler(
			rabbitmq.NewDefaultHandler(
				func(context.Context, interface{}) error {
					return nil
				},
			),
		),
	)

	s.assert.Error(err)
	s.assert.EqualError(err, "queue and exchangee must not be empty")
}
