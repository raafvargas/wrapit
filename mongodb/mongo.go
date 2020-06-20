package mongodb

import (
	"context"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	mongotrace "go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver"
)

// MongoConfig ...
type MongoConfig struct {
	ConnectionString string `yaml:"connection_string"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Database         string `yaml:"database"`
}

// Connect ...
func Connect(ctx context.Context, serviceName string, config *MongoConfig) (*mongo.Client, error) {
	uri, err := url.Parse(config.ConnectionString)

	if err != nil {
		return nil, err
	}

	if config.Username != "" && config.Password != "" {
		uri.User = url.UserPassword(config.Username, config.Password)
	}

	opt := options.Client().
		ApplyURI(uri.String()).
		SetMonitor(
			mongotrace.NewMonitor(serviceName),
		).
		SetAppName(serviceName)

	client, err := mongo.NewClient(opt)

	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Nearest()); err != nil {
		return nil, err
	}

	return client, nil
}
