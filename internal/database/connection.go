package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectMongoDB initializes a connection to MongoDB
func ConnectMongoDB(uri string) error {
    // Create a context with a longer timeout for the connection
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    log.Printf("MongoDB connection string: %s", uri)

    // Connect to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if (err != nil) {
        log.Printf("Error connecting to MongoDB: %v", err)
        return fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Create a separate context for the Ping operation
    pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer pingCancel()

    // Verify the connection
    err = client.Ping(pingCtx, nil)
    if (err != nil) {
        log.Printf("Error pinging MongoDB: %v", err)
        return fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    MongoClient = client
    return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(databaseName, collectionName string) *mongo.Collection {
    if MongoClient == nil {
        panic("MongoClient is not initialized. Call ConnectMongoDB first.")
    }
    return MongoClient.Database(databaseName).Collection(collectionName)
}

// DisconnectMongoDB closes the MongoDB connection
func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v\n", err)
		} else {
			log.Println("Disconnected from MongoDB!")
		}
	}
}