package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raafvargas/wrapit/configuration"
	"github.com/raafvargas/wrapit/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestMongo(t *testing.T) {
	cfg := new(configuration.Config)
	err := configuration.FromYAML("../tests/config.yaml", cfg)

	if err != nil {
		t.Fatal(err)
	}

	_, err = mongodb.Connect(context.Background(), uuid.New().String(), cfg.Mongo)
	assert.NoError(t, err)
}

func TestMongoWithPwd(t *testing.T) {
	cfg := new(configuration.Config)
	err := configuration.FromYAML("../tests/config.yaml", cfg)

	if err != nil {
		t.Fatal(err)
	}

	cfg.Mongo.Username = uuid.New().String()
	cfg.Mongo.Password = uuid.New().String()

	_, err = mongodb.Connect(context.Background(), uuid.New().String(), cfg.Mongo)
	assert.Error(t, err)
}

func TestMongoInvalidUri(t *testing.T) {
	_, err := mongodb.Connect(context.Background(), uuid.New().String(), &mongodb.MongoConfig{
		ConnectionString: "invalid_uri",
	})
	assert.Error(t, err)
}

func TestMongoInvalidServer(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()

	_, err := mongodb.Connect(ctx, uuid.New().String(), &mongodb.MongoConfig{
		ConnectionString: "mongodb://invalid_server:27017",
	})
	assert.Error(t, err)
}
