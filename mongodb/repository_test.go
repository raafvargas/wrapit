package mongodb_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/raafvargas/wrapit/configuration"
	"github.com/raafvargas/wrapit/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryTestSuite struct {
	suite.Suite
	assert *assert.Assertions

	repository *mongodb.MongoRepository
	config     *configuration.Config
	client     *mongo.Client
	collection *mongo.Collection
}

type MongoTestDocument struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Value string             `bson:"value"`
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.config = new(configuration.Config)

	err := configuration.FromYAML("../tests/config.yaml", s.config)
	s.assert.NoError(err)

	s.client, _ = mongodb.Connect(context.Background(), "", s.config.Mongo)
	s.collection = s.client.Database(s.config.Mongo.Database).Collection("tests")
	s.repository = mongodb.NewMongoRepository("ID", reflect.TypeOf(MongoTestDocument{}), s.collection)
}

func (s *RepositoryTestSuite) TestRepository() {
	doc := &MongoTestDocument{
		Value: uuid.New().String(),
	}

	err := s.repository.Insert(context.Background(), doc)

	s.assert.NoError(err)
	s.assert.False(doc.ID.IsZero())

	oldValue := doc.Value
	doc.Value = uuid.New().String()

	err = s.repository.Update(context.Background(), doc.ID, doc)
	s.assert.NoError(err)

	result2, err := s.repository.FindByID(context.Background(), doc.ID)
	s.assert.NoError(err)
	s.assert.NotEqual(oldValue, result2.(*MongoTestDocument).Value)

	err = s.repository.Delete(context.Background(), doc.ID)
	s.assert.NoError(err)

	_, err = s.repository.FindByID(context.Background(), doc.ID)
	s.assert.Error(mongo.ErrNoDocuments)
}

func (s *RepositoryTestSuite) TestRepositoryInvalidID() {
	doc := &struct {
		Value string `bson:"value"`
	}{
		Value: uuid.New().String(),
	}

	err := s.repository.Insert(context.Background(), doc)

	s.assert.Error(err)
	s.assert.EqualError(err, mongodb.ErrCannotSetID.Error())
}
