package database

import (
	"context"
	"errors"
	"fmt"
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

func StartSession() (mongo.Session, error) {
	if Client == nil {
		return nil, errors.New("No connection")
	}
	return Client.client.StartSession()
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

func (coll *MongoCollection) FindAllWithTimeout(filter interface{}, timeout time.Duration) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return coll.collection.Find(ctx, filter)
}

func (coll *MongoCollection) FindAll(filter interface{}) (*mongo.Cursor, error) {
	return coll.FindAllWithTimeout(filter, thirtySeconds)
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

func (coll *MongoCollection) UpdateByIDWithTimeout(ID string, update interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	filter := bson.M{idField: ObjectID}
	fmt.Println(ObjectID)
	opts := options.Update().SetUpsert(false)
	_, err = coll.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (coll *MongoCollection) UpdateByID(ID string, update interface{}) error {
	return coll.UpdateByIDWithTimeout(ID, update, tenSeconds)
}
