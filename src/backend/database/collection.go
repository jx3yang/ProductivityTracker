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

func (coll *MongoCollection) FindByIDWithTimeout(ID string, timeout time.Duration, ctx context.Context) (*mongo.SingleResult, error) {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	res := coll.collection.FindOne(ctx, bson.M{constants.IDField: ObjectID})
	return res, res.Err()
}

func (coll *MongoCollection) FindByID(ID string, ctx context.Context) (*mongo.SingleResult, error) {
	return coll.FindByIDWithTimeout(ID, tenSeconds, ctx)
}

func (coll *MongoCollection) FindAllWithTimeout(filter interface{}, timeout time.Duration, ctx context.Context) (*mongo.Cursor, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	return coll.collection.Find(ctx, filter)
}

func (coll *MongoCollection) FindAll(filter interface{}, ctx context.Context) (*mongo.Cursor, error) {
	return coll.FindAllWithTimeout(filter, thirtySeconds, ctx)
}

func (coll *MongoCollection) InsertOneWithTimeout(document interface{}, timeout time.Duration, ctx context.Context) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	res, err := coll.collection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (coll *MongoCollection) InsertOne(document interface{}, ctx context.Context) (string, error) {
	return coll.InsertOneWithTimeout(document, tenSeconds, ctx)
}

func (coll *MongoCollection) UpdateByIDWithTimeout(ID string, update interface{}, timeout time.Duration, ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	filter := bson.M{constants.IDField: ObjectID}
	opts := options.Update().SetUpsert(false)
	_, err = coll.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (coll *MongoCollection) UpdateByID(ID string, update interface{}, ctx context.Context) error {
	return coll.UpdateByIDWithTimeout(ID, update, tenSeconds, ctx)
}

func (coll *MongoCollection) BulkUpdateByIDsWithTimeout(idsToUpdate map[string]interface{}, timeout time.Duration, ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	var models []mongo.WriteModel

	for id, update := range idsToUpdate {
		ObjectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		filter := bson.M{constants.IDField: ObjectID}
		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
		models = append(models, updateModel)
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := coll.collection.BulkWrite(ctx, models, opts)
	return err
}

func (coll *MongoCollection) BulkUpdateByIDs(idsToUpdate map[string]interface{}, ctx context.Context) error {
	return coll.BulkUpdateByIDsWithTimeout(idsToUpdate, thirtySeconds, ctx)
}
