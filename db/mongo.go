// Package db mongo.go
package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Connect establishes a connection to the MongoDB database.
func Connect() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Update the URI as needed.
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// GetUsersCollection returns a handle to the users collection in the database.
func GetUsersCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("mySplit").Collection("users")
}

func GetGroupsCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("mySplit").Collection("groups")
}

func GetExpenseCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("mySplit").Collection("expenses")
}
