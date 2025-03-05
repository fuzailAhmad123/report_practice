package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a new client of Redis database.
func NewRedisClient() (*RedisClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: os.Getenv("REDIS_USERNAME"),
		DB:       0, // use default DB
	})
	// Test the connection to Redis
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return nil, err
	}
	fmt.Println("Connected to Redis")
	return &RedisClient{Client: client}, nil
}

// Function to close redis client.
func (rc *RedisClient) Close() {
	err := rc.Client.Close()
	if err != nil {
		fmt.Println("Failed to disconnect from Redis:", err)
		return
	}
	fmt.Println("Disconnected from Redis")
}
