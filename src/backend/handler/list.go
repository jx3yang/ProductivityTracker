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
