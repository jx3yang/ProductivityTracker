package handler

import (
	"context"

	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	model "github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var listCollection *db.MongoCollection

func initListCollection(d *db.MongoDatabase) {
	listCollection = d.InitCollection(listCollectionName)
}

func FindListByID(ID string) (*model.List, error) {
	res, err := listCollection.FindByID(ID)
	if err != nil {
		return nil, err
	}
	list := List{}
	res.Decode(&list)
	cards, err := FindAllFromList(ID)
	if err != nil {
		return nil, err
	}

	idToCardMap := make(map[string]*model.Card, len(cards))
	for idx := range cards {
		card := cards[idx]
		idToCardMap[card.ID] = card
	}

	var orderedCards []*model.CardMetaData
	for _, id := range list.CardOrder {
		card, ok := idToCardMap[id]
		if !ok {
			panic("Could not find id " + id + " in the fetched cards")
		}
		orderedCards = append(orderedCards,
			&model.CardMetaData{
				ID:      card.ID,
				Name:    card.Name,
				DueDate: card.DueDate,
			},
		)
	}

	return &model.List{
		ID:            list.ID.Hex(),
		Name:          list.Name,
		ParentBoardID: list.ParentBoardID,
		Cards:         orderedCards,
	}, nil
}

func CreateList(list *model.NewList) (*model.List, error) {
	res, err := boardCollection.FindByID(list.ParentBoardID)
	if err != nil {
		return nil, err
	}

	var board Board
	res.Decode(&board)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		res, err := listCollection.InsertOne(list)
		if err != nil {
			return nil, err
		}
		newOrder := append(board.ListOrder, res)
		update := bson.M{"$set": bson.M{"listOrder": newOrder}}
		err = boardCollection.UpdateByID(board.ID.Hex(), update)
		if err != nil {
			return nil, err
		}

		return &model.List{
			ID:    res,
			Name:  list.Name,
			Cards: nil,
		}, nil
	}

	session, err := db.StartSession()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer session.EndSession(ctx)
	result, err := session.WithTransaction(ctx, callback)
	return result.(*model.List), err
}
