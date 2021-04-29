package handler

import (
	"errors"

	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	model "github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
)

var boardCollection *db.MongoCollection

func initBoardCollection(d *db.MongoDatabase) {
	boardCollection = d.InitCollection(boardCollectionName)
}

func FindBoardByID(ID string) (*model.Board, error) {
	res, err := boardCollection.FindByID(ID, nil)
	if err != nil {
		return nil, err
	}

	board := db.Board{}
	res.Decode(&board)
	lists, err := FindAllUnarchivedListsFromBoard(ID)
	if err != nil {
		return nil, err
	}

	idToListsMap := make(map[string]*model.List, len(lists))
	for idx := range lists {
		list := lists[idx]
		idToListsMap[list.ID] = list
	}

	var orderedLists []*model.List
	for _, id := range board.ListOrder {
		list, ok := idToListsMap[id]
		if !ok {
			return nil, errors.New("Could not find id " + id + " in the fetched lists")
		}
		orderedLists = append(orderedLists, list)
	}

	return &model.Board{
		ID:    board.ID.Hex(),
		Name:  board.Name,
		Lists: orderedLists,
	}, nil
}

func CreateBoard(newBoard *model.NewBoard) (*model.Board, error) {
	document := map[string]interface{}{
		"name":      newBoard.Name,
		"listOrder": make([]interface{}, 0),
	}

	res, err := boardCollection.InsertOne(document, nil)
	if err != nil {
		return nil, err
	}
	return &model.Board{
		ID:    res,
		Name:  newBoard.Name,
		Lists: make([]*model.List, 0),
	}, nil
}
