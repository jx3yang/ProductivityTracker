package handler

import (
	"context"
	"errors"

	"github.com/jx3yang/ProductivityTracker/src/backend/constants"
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	model "github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var listCollection *db.MongoCollection

func initListCollection(d *db.MongoDatabase) {
	listCollection = d.InitCollection(listCollectionName)
}

func FindAllUnarchivedListsFromBoard(boardID string) ([]*model.List, error) {
	cursor, err := listCollection.FindAll(
		bson.M{constants.ParentBoardIDField: boardID, constants.ArchivedField: bson.M{"$ne": true}}, nil,
	)
	if err != nil {
		return nil, err
	}
	var dbLists []*db.List
	if err = cursor.All(context.TODO(), &dbLists); err != nil {
		return nil, err
	}

	var lists []*model.List
	for _, dbList := range dbLists {
		list, err := FindListByID(dbList.ID.Hex())
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	return lists, nil
}

func FindListByID(ID string) (*model.List, error) {
	res, err := listCollection.FindByID(ID, nil)
	if err != nil {
		return nil, err
	}
	list := db.List{}
	res.Decode(&list)
	cards, err := FindAllUnarchivedCardsFromList(ID)
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
			return nil, errors.New("Could not find id " + id + " in the fetched cards")
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
	operation := func(sessCtx mongo.SessionContext) (interface{}, error) {
		res, err := boardCollection.FindByID(list.ParentBoardID, sessCtx)
		if err != nil {
			return nil, err
		}

		var board db.Board
		res.Decode(&board)

		document := map[string]interface{}{
			"name":          list.Name,
			"parentBoardId": list.ParentBoardID,
			"cardOrder":     make([]interface{}, 0),
		}

		listID, err := listCollection.InsertOne(document, sessCtx)
		if err != nil {
			return nil, err
		}
		newOrder := append(board.ListOrder, listID)
		update := bson.M{"$set": bson.M{constants.ListOrderField: newOrder}}
		err = boardCollection.UpdateByID(board.ID.Hex(), update, sessCtx)
		if err != nil {
			return nil, err
		}

		return &model.List{
			ID:    listID,
			Name:  list.Name,
			Cards: make([]*model.CardMetaData, 0),
		}, nil
	}

	result, err := executeWithSession(operation)
	return result.(*model.List), err
}

func UpdateListOrder(changeListOrder *model.ChangeListOrder) (bool, error) {
	if changeListOrder.SrcIdx == changeListOrder.DestIdx {
		return true, nil
	}

	operation := func(sessCtx mongo.SessionContext) (interface{}, error) {
		boardID := changeListOrder.BoardID
		listID := changeListOrder.ListID
		srcIdx := changeListOrder.SrcIdx

		res, err := boardCollection.FindByID(boardID, sessCtx)
		if err != nil {
			return false, err
		}

		board := db.Board{}
		res.Decode(&board)

		listOrder := board.ListOrder

		if srcIdx >= len(listOrder) || listOrder[srcIdx] != listID {
			return false, errors.New("The board state is modified")
		}

		destIdx := changeListOrder.DestIdx
		if len(listOrder) <= destIdx {
			destIdx = len(listOrder) - 1
		}

		if srcIdx == destIdx {
			return true, nil
		}

		newOrder := moveElement(listOrder, srcIdx, destIdx)

		update := bson.M{"$set": bson.M{constants.ListOrderField: newOrder}}
		err = boardCollection.UpdateByID(board.ID.Hex(), update, sessCtx)

		return err == nil, err
	}

	result, err := executeWithSession(operation)
	return result.(bool), err
}

func ArchiveList(list *model.ListIdentifier) (bool, error) {
	return true, nil
}

func DeleteList(list *model.ListIdentifier) (bool, error) {
	return true, nil
}
