package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Function to do mongo aggregation.
func Aggregation(collection *mongo.Collection, matchFields bson.D, groupByFields bson.D, projectFields bson.D) ([]interface{}, error) {
	// Match Stage
	matchStage := bson.D{{Key: "$match", Value: matchFields}}

	// Group Stage
	groupStage := bson.D{{Key: "$group", Value: groupByFields}}

	// Projection Stage
	projectStage := bson.D{{Key: "$project", Value: projectFields}}

	pipeline := mongo.Pipeline{matchStage, groupStage, projectStage}
	options := options.Aggregate().SetAllowDiskUse(true).SetMaxTime(30 * time.Second)
	cur, err := collection.Aggregate(context.Background(), pipeline, options)
	var result []interface{}
	if err != nil {
		return result, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var doc interface{}
		err := cur.Decode(&doc)
		if err != nil {
			return result, err
		}
		result = append(result, doc)
	}

	return result, nil
}
