package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
}

type MongoDefaultDatabase struct {
	Db *mongo.Database
}

func ConnectWithMongoDb(connectionStr string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionStr))
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB")
	return &MongoClient{Client: client}, nil
}

// Function to close mongodb client.
func (mc *MongoClient) Close() {
	err := mc.Client.Disconnect(context.Background())
	if err != nil {
		fmt.Println("Failed to disconnect from MongoDB:", err)
		return
	}
	fmt.Println("Disconnected from MongoDB")
}

// Function to create a new mongodb database.
func (mc *MongoClient) NewDatabase(dbName string) *MongoDefaultDatabase {
	return &MongoDefaultDatabase{Db: mc.Client.Database(dbName)}
}
