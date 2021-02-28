package handler

import (
	"context"
	"errors"
	"fmt"

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

func updateCardOrderSameList(changeCardOrder *model.ChangeCardOrder) (bool, error) {
	srcListID := changeCardOrder.SrcListID
	srcIdx := changeCardOrder.SrcIdx
	cardID := changeCardOrder.CardID

	res, err := listCollection.FindByID(srcListID)
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
	err = listCollection.UpdateByID(srcListID, update)

	return err == nil, err
}

func updateCardOrderDifferentLists(changeCardOrder *model.ChangeCardOrder) (bool, error) {
	srcListID := changeCardOrder.SrcListID
	destListID := changeCardOrder.DestListID
	srcIdx := changeCardOrder.SrcIdx
	cardID := changeCardOrder.CardID

	res, err := listCollection.FindByID(srcListID)
	if err != nil {
		return false, err
	}

	srcList := db.List{}
	res.Decode(&srcList)

	if srcList.ParentBoardID != changeCardOrder.BoardID {
		return false, errors.New("Source list does not belong to board with id " + changeCardOrder.BoardID)
	}

	res, err = listCollection.FindByID(destListID)
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
	if len(destListCardOrder) <= destIdx {
		destIdx = len(destListCardOrder) - 1
	}

	newSrcOrder := removeOneFromList(srcListCardOrder, srcIdx)
	newDestOrder := addOneToList(destListCardOrder, destIdx, cardID)

	fmt.Println(newSrcOrder)
	fmt.Println(newDestOrder)

	idsToUpdate := make(map[string]interface{}, 2)
	idsToUpdate[srcListID] = bson.M{"$set": bson.M{constants.CardOrderField: newSrcOrder}}
	idsToUpdate[destListID] = bson.M{"$set": bson.M{constants.CardOrderField: newDestOrder}}

	err = listCollection.BulkUpdateByIDs(idsToUpdate)

	if err != nil {
		return false, err
	}

	res, err = cardCollection.FindByID(cardID)
	if err != nil {
		return false, err
	}

	card := db.Card{}
	res.Decode(&card)

	parentListUpdate := bson.M{"$set": bson.M{constants.ParentListIDField: destListID}}
	err = cardCollection.UpdateByID(cardID, parentListUpdate)

	return err == nil, err
}

func UpdateCardOrder(changeCardOrder *model.ChangeCardOrder) (bool, error) {
	if changeCardOrder.SrcListID == changeCardOrder.DestListID && changeCardOrder.SrcIdx == changeCardOrder.DestIdx {
		return true, nil
	}

	operation := func() (interface{}, error) {
		if changeCardOrder.SrcListID == changeCardOrder.DestListID {
			return updateCardOrderSameList(changeCardOrder)
		}
		return updateCardOrderDifferentLists(changeCardOrder)
	}

	result, err := executeWithSession(operation)
	return result.(bool), err
}
