package database

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

func StartSession() (mongo.Session, error) {
	if Client == nil {
		return nil, errors.New("No connection")
	}
	return Client.client.StartSession()
}
