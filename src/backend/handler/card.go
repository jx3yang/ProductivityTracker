package handler

import (
	"context"
	"fmt"

	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	"github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func FindAllFromList(listID string) ([]*model.Card, error) {
	cursor, err := cardCollection.FindAll(bson.M{"parentListId": listID})
	if err != nil {
		return nil, err
	}
	var cards []*model.Card
	if err = cursor.All(context.TODO(), &cards); err != nil {
		return nil, err
	}
	fmt.Println(cards)
	return cards, nil
}

func CreateCard(card *model.NewCard) (*model.Card, error) {
	res, err := listCollection.FindByID(card.ParentListID)
	if err != nil {
		return nil, err
	}
	list := List{}
	res.Decode(&list)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		res, err := cardCollection.InsertOne(card)
		if err != nil {
			return nil, err
		}
		newOrder := append(list.CardOrder, res)
		update := bson.M{"$set": bson.M{"cardOrder": newOrder}}
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

	session, err := db.StartSession()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer session.EndSession(ctx)
	result, err := session.WithTransaction(ctx, callback)
	return result.(*model.Card), err
}
