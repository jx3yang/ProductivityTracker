package handler

import (
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	model "github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
)

var boardCollection *db.MongoCollection

func initBoardCollection(d *db.MongoDatabase) {
	boardCollection = d.InitCollection(boardCollectionName)
}

func FindBoardByID(ID string) (*model.Board, error) {
	res, err := boardCollection.FindByID(ID)
	if err != nil {
		return nil, err
	}

	board := Board{}
	res.Decode(&board)
	lists, err := FindAllListsFromBoard(ID)
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
			panic("Could not find id " + id + " in the fetched lists")
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
	res, err := boardCollection.InsertOne(newBoard)
	if err != nil {
		return nil, err
	}
	return &model.Board{
		ID:    res,
		Name:  newBoard.Name,
		Lists: nil,
	}, nil
}
