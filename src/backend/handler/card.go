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

func CreateCard(card *model.NewCard) (*model.Card, error) {
	_, err := FindListByID(card.ParentListID)
	if err != nil {
		return nil, err
	}
	res, err := cardCollection.InsertOne(card)
	if err != nil {
		return nil, err
	}
	return &model.Card{
		ID:           res,
		Name:         card.Name,
		DueDate:      card.DueDate,
		ParentListID: card.ParentListID,
	}, nil
}
