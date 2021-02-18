package handler

import (
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	model "github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
)

var cardCollection *db.MongoCollection

func initCardCollection(d *db.MongoDatabase) {
	cardCollection = d.InitCollection(cardCollectionName)
}

func FindCardByID(ID string) (*model.Card, error) {
	res, err := cardCollection.FindByID(ID)
	if err != nil {
		return nil, err
	}
	card := model.Card{}
	res.Decode(&card)
	return &card, nil
}
