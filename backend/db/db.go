package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(options *options.ClientOptions) *mongo.Client{
	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		panic(err)
	}
	return client
}

func NewMongoDatabase(client *mongo.Client, databaseName string) *mongo.Database {
	return client.Database(databaseName)
}

func NewMongoCollection(database *mongo.Database, collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}