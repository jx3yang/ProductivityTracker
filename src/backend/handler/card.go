package handler

import (
	"context"
	"errors"

	"github.com/jx3yang/ProductivityTracker/src/backend/constants"
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	"github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	res, err := cardCollection.FindByID(ID, nil)
	if err != nil {
		return nil, err
	}
	card := db.Card{}
	res.Decode(&card)

	return dbCardToGraphqlCard(&card), nil
}

func FindAllUnarchivedCardsFromList(listID string) ([]*model.Card, error) {
	cursor, err := cardCollection.FindAll(
		bson.M{constants.ParentListIDField: listID, constants.ArchivedField: bson.M{"$ne": true}}, nil,
	)
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
	operation := func(sessCtx mongo.SessionContext) (interface{}, error) {
		res, err := listCollection.FindByID(card.ParentListID, sessCtx)
		if err != nil {
			return nil, err
		}
		list := db.List{}
		res.Decode(&list)

		if list.ParentBoardID != card.ParentBoardID {
			return nil, errors.New("The list with id " + card.ParentListID + " does not belong to board " + card.ParentBoardID)
		}

		cardID, err := cardCollection.InsertOne(card, sessCtx)
		if err != nil {
			return nil, err
		}
		newOrder := append(list.CardOrder, cardID)
		update := bson.M{"$set": bson.M{constants.CardOrderField: newOrder}}
		err = listCollection.UpdateByID(list.ID.Hex(), update, sessCtx)
		if err != nil {
			return nil, err
		}
		return &model.Card{
			ID:            cardID,
			Name:          card.Name,
			DueDate:       card.DueDate,
			ParentListID:  card.ParentListID,
			ParentBoardID: card.ParentBoardID,
		}, nil
	}

	result, err := executeWithSession(operation)
	if result == nil {
		return nil, err
	}
	return result.(*model.Card), err
}

func updateCardOrderSameList(changeCardOrder *model.ChangeCardOrder, ctx context.Context) (bool, error) {
	srcListID := changeCardOrder.SrcListID
	srcIdx := changeCardOrder.SrcIdx
	cardID := changeCardOrder.CardID

	res, err := listCollection.FindByID(srcListID, ctx)
	if err != nil {
		return false, err
	}

	lst := db.List{}
	res.Decode(&lst)

	if lst.ParentBoardID != changeCardOrder.BoardID {
		return false, errors.New("The list does not belong to board with id " + changeCardOrder.BoardID)
	}

	cardOrder := lst.CardOrder

	if srcIdx >= len(cardOrder) || cardOrder[srcIdx] != cardID {
		return false, errors.New("The list state is modified")
	}

	destIdx := changeCardOrder.DestIdx
	if len(cardOrder) <= destIdx {
		destIdx = len(cardOrder) - 1
	}

	if srcIdx == destIdx {
		return true, nil
	}

	newOrder := moveElement(cardOrder, srcIdx, destIdx)

	update := bson.M{"$set": bson.M{constants.CardOrderField: newOrder}}
	err = listCollection.UpdateByID(srcListID, update, ctx)

	return err == nil, err
}

func updateCardOrderDifferentLists(changeCardOrder *model.ChangeCardOrder, ctx context.Context) (bool, error) {
	srcListID := changeCardOrder.SrcListID
	destListID := changeCardOrder.DestListID
	srcIdx := changeCardOrder.SrcIdx
	cardID := changeCardOrder.CardID

	res, err := listCollection.FindByID(srcListID, ctx)
	if err != nil {
		return false, err
	}

	srcList := db.List{}
	res.Decode(&srcList)

	if srcList.ParentBoardID != changeCardOrder.BoardID {
		return false, errors.New("Source list does not belong to board with id " + changeCardOrder.BoardID)
	}

	res, err = listCollection.FindByID(destListID, ctx)
	if err != nil {
		return false, err
	}

	destList := db.List{}
	res.Decode(&destList)

	if destList.ParentBoardID != changeCardOrder.BoardID {
		return false, errors.New("Destination list does not belong to board with id " + changeCardOrder.BoardID)
	}

	srcListCardOrder := srcList.CardOrder
	destListCardOrder := destList.CardOrder

	if srcIdx >= len(srcListCardOrder) || srcListCardOrder[srcIdx] != cardID {
		return false, errors.New("The list state is modified")
	}

	destIdx := changeCardOrder.DestIdx
	if len(destListCardOrder) < destIdx {
		destIdx = len(destListCardOrder)
	}

	newSrcOrder := removeOneFromList(srcListCardOrder, srcIdx)
	newDestOrder := addOneToList(destListCardOrder, destIdx, cardID)

	idsToUpdate := make(map[string]interface{}, 2)
	idsToUpdate[srcListID] = bson.M{"$set": bson.M{constants.CardOrderField: newSrcOrder}}
	idsToUpdate[destListID] = bson.M{"$set": bson.M{constants.CardOrderField: newDestOrder}}

	err = listCollection.BulkUpdateByIDs(idsToUpdate, ctx)

	if err != nil {
		return false, err
	}

	res, err = cardCollection.FindByID(cardID, ctx)
	if err != nil {
		return false, err
	}

	card := db.Card{}
	res.Decode(&card)

	parentListUpdate := bson.M{"$set": bson.M{constants.ParentListIDField: destListID}}
	err = cardCollection.UpdateByID(cardID, parentListUpdate, ctx)

	return err == nil, err
}

func UpdateCardOrder(changeCardOrder *model.ChangeCardOrder) (bool, error) {
	if changeCardOrder.SrcListID == changeCardOrder.DestListID && changeCardOrder.SrcIdx == changeCardOrder.DestIdx {
		return true, nil
	}

	operation := func(sessCtx mongo.SessionContext) (interface{}, error) {
		if changeCardOrder.SrcListID == changeCardOrder.DestListID {
			return updateCardOrderSameList(changeCardOrder, sessCtx)
		}
		return updateCardOrderDifferentLists(changeCardOrder, sessCtx)
	}

	result, err := executeWithSession(operation)
	return result.(bool), err
}

func ArchiveCard(card *model.CardIdentifier) (bool, error) {
	operation := func(sessCtx mongo.SessionContext) (interface{}, error) {
		res, err := listCollection.FindByID(card.ParentListID, sessCtx)
		if err != nil {
			return nil, err
		}
		list := db.List{}
		res.Decode(&list)

		if list.ParentBoardID != card.ParentBoardID {
			return nil, errors.New("The list with id " + card.ParentListID + " does not belong to board " + card.ParentBoardID)
		}

		targetIdx := -1
		for idx, cardID := range list.CardOrder {
			if cardID == card.ID {
				targetIdx = idx
				break
			}
		}
		if targetIdx == -1 {
			return nil, errors.New("Card with id " + card.ID + " does not belong to list " + card.ParentListID)
		}

		idx := 0
		newOrder := make([]string, len(list.CardOrder)-1)
		for i := 0; i < len(newOrder); i++ {
			if idx == targetIdx {
				idx++
			}
			newOrder[i] = list.CardOrder[idx]
			idx++
		}

		listUpdate := bson.M{"$set": bson.M{constants.CardOrderField: newOrder}}
		err = listCollection.UpdateByID(card.ParentListID, listUpdate, sessCtx)
		if err != nil {
			return nil, err
		}

		cardUpdate := bson.M{"$set": bson.M{constants.ArchivedField: true}}
		err = cardCollection.UpdateByID(card.ID, cardUpdate, sessCtx)

		return err == nil, err
	}

	result, err := executeWithSession(operation)
	return result.(bool), err
}

func DeleteCard(card *model.CardIdentifier) (bool, error) {
	return true, nil
}
