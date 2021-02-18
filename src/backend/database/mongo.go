package database

import (
	"context"
	"time"

	"github.com/jx3yang/ProductivityTracker/src/backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tenSeconds = 10 * time.Second
const thirtySeconds = 30 * time.Second

const idField = "_id"

func ConnectionFromConfig(c *config.Config) (*MongoConnection, error) {
	uri := "mongodb://" + c.DBUsername + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort
	return newMongoConnection(uri)
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

func (coll *MongoCollection) FindByIDWithTimeout(ID string, timeout time.Duration) (*mongo.SingleResult, error) {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return coll.collection.FindOne(ctx, bson.M{idField: ObjectID}), nil
}

func (coll *MongoCollection) FindByID(ID string) (*mongo.SingleResult, error) {
	return coll.FindByIDWithTimeout(ID, tenSeconds)
}

func (coll *MongoCollection) InsertOneWithTimeout(document interface{}, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := coll.collection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (coll *MongoCollection) InsertOne(document interface{}) (string, error) {
	return coll.InsertOneWithTimeout(document, tenSeconds)
}
