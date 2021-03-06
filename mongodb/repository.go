package mongodb

import (
	"context"
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrCannotSetID ...
	ErrCannotSetID = errors.New("cannot set id property")
)

// Repository ...
type Repository interface {
	Insert(ctx context.Context, document interface{}) error
	Update(ctx context.Context, id interface{}, document interface{}) error
	FindByID(ctx context.Context, id interface{}) (interface{}, error)
}

// MongoRepository ...
type MongoRepository struct {
	idProperty   string
	documentType reflect.Type

	Collection *mongo.Collection
}

// NewMongoRepository ...
func NewMongoRepository(idProperty string, documentType reflect.Type, collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		idProperty:   idProperty,
		documentType: documentType,
		Collection:   collection,
	}
}

// Insert ...
func (r *MongoRepository) Insert(ctx context.Context, document interface{}) error {
	result, err := r.Collection.InsertOne(ctx, document)

	if err != nil {
		return err
	}

	field := reflect.ValueOf(document).Elem().FieldByName(r.idProperty)

	if !field.IsValid() || !field.CanSet() {
		logrus.WithField("document", document).
			Warnf("property %s is invalid or cannot be set", r.idProperty)
		return ErrCannotSetID
	}

	field.Set(reflect.ValueOf(result.InsertedID))

	return err
}

// Update ...
func (r *MongoRepository) Update(ctx context.Context, id interface{}, document interface{}) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{"$set": document})

	return err
}

// FindByID ...
func (r *MongoRepository) FindByID(ctx context.Context, id interface{}) (interface{}, error) {
	doc := reflect.New(r.documentType).Interface()

	err := r.Collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(doc)

	return doc, err
}

// Delete ...
func (r *MongoRepository) Delete(ctx context.Context, id interface{}) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})

	return err
}
