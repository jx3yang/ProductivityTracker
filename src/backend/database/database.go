package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnection interface {
	FindByID(ID string, collection string) (*mongo.SingleResult, error)
	FindByIDWithTimeout(ID string, collection string, timeout time.Duration) (*mongo.SingleResult, error)
	InsertOne(document interface{}, collection string) (string, error)
	InsertOneWithTimeout(document interface{}, collection string, timeout time.Duration) (string, error)
}

// mongoConnection holds a connection to a MongoDB
type mongoConnection struct {
	dbName string
	client *mongo.Client
}

const tenSeconds = 10 * time.Second
const thirtySeconds = 30 * time.Second

// NewConnection returns a new DB connection
func NewConnection(uri string, dbName string) (DBConnection, error) {
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

	return &mongoConnection{
		dbName: dbName,
		client: client,
	}, nil
}

func (db *mongoConnection) getCollection(collection string) *mongo.Collection {
	client := db.client
	dbName := db.dbName
	return client.Database(dbName).Collection(collection)
}

func (db *mongoConnection) FindByIDWithTimeout(ID string, collection string, timeout time.Duration) (*mongo.SingleResult, error) {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	dbCollection := db.getCollection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return dbCollection.FindOne(ctx, bson.M{"_id": ObjectID}), nil
}

func (db *mongoConnection) FindByID(ID string, collection string) (*mongo.SingleResult, error) {
	return db.FindByIDWithTimeout(ID, collection, tenSeconds)
}

func (db *mongoConnection) InsertOneWithTimeout(document interface{}, collection string, timeout time.Duration) (string, error) {
	dbCollection := db.getCollection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := dbCollection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *mongoConnection) InsertOne(document interface{}, collection string) (string, error) {
	return db.InsertOneWithTimeout(document, collection, tenSeconds)
}
