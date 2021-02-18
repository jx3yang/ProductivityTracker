package handler

import (
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	model "github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
)

const cardCollectionName = "card"

var cardCollection *db.MongoCollection

func InitCardCollection(m *db.MongoConnection) {
	cardCollection = m.InitCollection(cardCollectionName)
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
