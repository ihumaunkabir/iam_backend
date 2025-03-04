package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database holds the MongoDB client and database connection
type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoConnection establishes a connection to MongoDB
func NewMongoConnection(uri, dbName string) (*Database, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri).SetServerSelectionTimeout(10 * time.Second)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")

	return &Database{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Disconnect closes the MongoDB connection
func (db *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.Client.Disconnect(ctx)
}
