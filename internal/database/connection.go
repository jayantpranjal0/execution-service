package database

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoClient *mongo.Client

// ConnectMongoDB initializes a connection to MongoDB
func ConnectMongoDB(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")
	MongoClient = client
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(databaseName, collectionName string) *mongo.Collection {
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