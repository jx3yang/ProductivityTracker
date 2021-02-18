package database

import "go.mongodb.org/mongo-driver/mongo"

type MongoConnection struct {
	dbName string
	client *mongo.Client
}

type MongoCollection struct {
	collection *mongo.Collection
}
