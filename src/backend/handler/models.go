package handler

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name"`
	ListOrder []string           `json:"listOrder"`
}

type List struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          string             `json:"name"`
	ParentBoardID string             `json:"parentBoardId"`
	CardOrder     []string           `json:"cardOrder"`
}
