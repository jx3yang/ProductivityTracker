package handler

import (
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	model "github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
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
	list := model.List{}
	res.Decode(&list)
	return &list, nil
}

func CreateList(list *model.NewList) (*model.List, error) {
	_, err := boardCollection.FindByID(list.ParentBoardID)
	if err != nil {
		return nil, err
	}
	id, err := listCollection.InsertOne(list)
	if err != nil {
		return nil, err
	}
	return &model.List{
		ID:   id,
		Name: list.Name,
	}, nil
}
