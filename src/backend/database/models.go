package database

import "go.mongodb.org/mongo-driver/mongo"

type MongoConnection struct {
	client *mongo.Client
}

type MongoDatabase struct {
	db *mongo.Database
}

type MongoCollection struct {
	collection *mongo.Collection
}
