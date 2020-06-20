package rabbitmq_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/raafvargas/wrapit/configuration"
	"github.com/raafvargas/wrapit/rabbitmq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProducerTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	config     *configuration.Config
	connection *rabbitmq.RabbitConnection
}

func TestProducerTestSuite(t *testing.T) {
	suite.Run(t, new(ProducerTestSuite))
}

func (s *ProducerTestSuite) SetupTest() {
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

	s.connection = connection
}

func (s *ProducerTestSuite) TestProducer() {
	queueName := uuid.New().String()
	exchangeName := uuid.New().String()

	err := s.connection.EnsureExchange(context.Background(), exchangeName)
	s.assert.NoError(err)
	err = s.connection.EnsureQueue(context.Background(), queueName, exchangeName)
	s.assert.NoError(err)

	producer := rabbitmq.NewProducer(s.connection)

	message := &struct {
		A string `json:"a"`
	}{
		A: "B",
	}

	err = producer.Publish(context.Background(), exchangeName, message)

	s.assert.NoError(err)
}

func (s *ProducerTestSuite) TestProducerFailJson() {
	exchangeName := uuid.New().String()

	err := s.connection.EnsureExchange(context.Background(), exchangeName)
	s.assert.NoError(err)

	producer := rabbitmq.NewProducer(s.connection)

	invalidArg := make(chan int, 1)
	defer close(invalidArg)

	err = producer.Publish(context.Background(), exchangeName, invalidArg)

	s.assert.Error(err)
}

func (s *ProducerTestSuite) TestProducerNoQueue() {
	exchangeName := uuid.New().String()

	err := s.connection.EnsureExchange(context.Background(), exchangeName)
	s.assert.NoError(err)

	producer := rabbitmq.NewProducer(s.connection)

	err = producer.Publish(context.Background(), exchangeName, "body")

	s.assert.NoError(err)
}
