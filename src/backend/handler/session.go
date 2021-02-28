package handler

import (
	"context"

	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func executeWithSession(operation func() (interface{}, error)) (interface{}, error) {
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return operation()
	}
	session, err := db.StartSession()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer session.EndSession(ctx)
	return session.WithTransaction(ctx, callback)
}
