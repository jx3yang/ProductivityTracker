package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	ListOrder []string           `json:"listOrder" bson:"listOrder"`
	Archived  bool               `json:"archived,omitempty" bson:"archived,omitempty"`
}

type List struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	ParentBoardID string             `json:"parentBoardId" bson:"parentBoardId"`
	CardOrder     []string           `json:"cardOrder" bson:"cardOrder"`
	Archived      bool               `json:"archived,omitempty" bson:"archived,omitempty"`
}

type Card struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	DueDate       *string            `json:"dueDate" bson:"dueDate"`
	ParentBoardID string             `json:"parentBoardId" bson:"parentBoardId"`
	ParentListID  string             `json:"parentListId" bson:"parentListId"`
	Archived      bool               `json:"archived,omitempty" bson:"archived,omitempty"`
}
