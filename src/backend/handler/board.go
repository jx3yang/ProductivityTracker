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
	board := model.Board{}
	res.Decode(&board)
	return &board, nil
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
