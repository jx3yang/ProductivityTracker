package database

import (
	"context"
	"time"

	"github.com/jx3yang/ProductivityTracker/src/backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tenSeconds = 10 * time.Second
const thirtySeconds = 30 * time.Second

var Client *MongoConnection

func InitConnectionFromConfig(c *config.Config) (*MongoConnection, error) {
	// TODO: add credentials
	uri := "mongodb://" + c.DBHost + ":" + c.DBPort
	client, err := newMongoConnection(uri)
	if err != nil {
		return nil, err
	}
	Client = client
	return Client, nil
}

func newMongoConnection(uri string) (*MongoConnection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), tenSeconds)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &MongoConnection{
		client: client,
	}, nil
}

func (conn *MongoConnection) InitDatabase(database string) *MongoDatabase {
	return &MongoDatabase{
		db: conn.client.Database(database),
	}
}

func (database *MongoDatabase) InitCollection(collection string) *MongoCollection {
	return &MongoCollection{
		collection: database.db.Collection(collection),
	}
}
