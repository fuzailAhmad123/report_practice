package model

import (
	"context"

	"github.com/fuzailAhmad123/test_report/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model interface {
	TableName() string
}

// InsertOne is a common db model function used to insert single document of interface T.
func InsertOne[T Model](ctx context.Context, db *mongodb.MongoDefaultDatabase, options *options.InsertOneOptions, doc *T) (*mongo.InsertOneResult, error) {
	var m T

	insertResult, err := db.Db.Collection(m.TableName()).InsertOne(ctx, doc, options)
	if err != nil {
		return nil, err
	}

	return insertResult, nil
}

// InsertOne is a common db model function used to insert multiple document of interface T.
func InsertMany[T Model](ctx context.Context, db *mongodb.MongoDefaultDatabase, options *options.InsertManyOptions, docs []T) (*mongo.InsertManyResult, error) {
	var m T

	var documents []interface{}
	for _, commission := range docs {
		documents = append(documents, commission)
	}

	res, err := db.Db.Collection(m.TableName()).InsertMany(ctx, documents, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindOne is the function to fetch document of type T that satisfies filter.
func FindOne[T Model](ctx context.Context, db *mongodb.MongoDefaultDatabase, filter bson.M, projection *options.FindOneOptions) (*T, error) {
	var m T
	err := db.Db.Collection(m.TableName()).FindOne(ctx, filter, projection).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &m, nil
}

// DeleteMany is the function that deletes all documents of type T which satisfies the filter.
func DeleteMany[T Model](ctx context.Context, db *mongodb.MongoDefaultDatabase, filter bson.M) error {
	var m T

	_, err := db.Db.Collection(m.TableName()).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
