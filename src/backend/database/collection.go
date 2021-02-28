package database

import (
	"context"
	"time"

	"github.com/jx3yang/ProductivityTracker/src/backend/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (coll *MongoCollection) FindByIDWithTimeout(ID string, timeout time.Duration) (*mongo.SingleResult, error) {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return coll.collection.FindOne(ctx, bson.M{constants.IDField: ObjectID}), nil
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
	filter := bson.M{constants.IDField: ObjectID}
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
