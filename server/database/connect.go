package mongoClient

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

func InitMongoClient() {
	clientOnce.Do(func() {
		//mongoHost := os.Getenv("MONGO_HOST")
		godotenv.Load()
		mongoPort := os.Getenv("MONGO_PORT")
		mongoHost := os.Getenv("MONGO_HOST")
		mongoURI := "mongodb://" + mongoHost + ":" + mongoPort
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		clientOptions := options.Client().ApplyURI(mongoURI)
		var err error
		clientInstance, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		err = clientInstance.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}
		log.Println("MongoDB connection established")
	})
}
func GetClient() *mongo.Client {
	return clientInstance
}
