package handler

import (
	"context"

	"github.com/jx3yang/ProductivityTracker/src/backend/constants"
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	"github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

var cardCollection *db.MongoCollection

func initCardCollection(d *db.MongoDatabase) {
	cardCollection = d.InitCollection(cardCollectionName)
}

func dbCardToGraphqlCard(card *db.Card) *model.Card {
	return &model.Card{
		ID:            card.ID.Hex(),
		Name:          card.Name,
		DueDate:       card.DueDate,
		ParentBoardID: card.ParentBoardID,
		ParentListID:  card.ParentListID,
	}
}

func FindCardByID(ID string) (*model.Card, error) {
	res, err := cardCollection.FindByID(ID)
	if err != nil {
		return nil, err
	}
	card := db.Card{}
	res.Decode(&card)

	return dbCardToGraphqlCard(&card), nil
}

func FindAllCardsFromList(listID string) ([]*model.Card, error) {
	cursor, err := cardCollection.FindAll(bson.M{constants.ParentListIDField: listID})
	if err != nil {
		return nil, err
	}
	var cards []*db.Card
	if err = cursor.All(context.TODO(), &cards); err != nil {
		return nil, err
	}

	var graphqlCards []*model.Card

	for _, card := range cards {
		graphqlCards = append(graphqlCards, dbCardToGraphqlCard(card))
	}
	return graphqlCards, nil
}

func CreateCard(card *model.NewCard) (*model.Card, error) {
	res, err := listCollection.FindByID(card.ParentListID)
	if err != nil {
		return nil, err
	}
	list := db.List{}
	res.Decode(&list)

	operation := func() (interface{}, error) {
		res, err := cardCollection.InsertOne(card)
		if err != nil {
			return nil, err
		}
		newOrder := append(list.CardOrder, res)
		update := bson.M{"$set": bson.M{constants.CardOrderField: newOrder}}
		err = listCollection.UpdateByID(list.ID.Hex(), update)
		if err != nil {
			return nil, err
		}
		return &model.Card{
			ID:            res,
			Name:          card.Name,
			DueDate:       card.DueDate,
			ParentListID:  card.ParentListID,
			ParentBoardID: card.ParentBoardID,
		}, nil
	}

	result, err := executeWithSession(operation)
	return result.(*model.Card), err
}
