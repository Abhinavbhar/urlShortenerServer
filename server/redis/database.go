package redis

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	client *redis.Client
	once   sync.Once
)

func InitRedis() {
	once.Do(func() {
		godotenv.Load()
		// Get Redis host and port from environment variables
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")
		redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

		client = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "", // No password set
			DB:       0,  // Use default DB
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := client.Ping(ctx).Err()
		if err != nil {
			fmt.Println("Could not ping Redis:", err)
		} else {
			fmt.Println("Redis ping successful")
		}
	})
}

func RedisDatabase() *redis.Client {
	return client
}
