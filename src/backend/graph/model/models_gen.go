// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Board struct {
	ID    string  `json:"_id" bson:"_id"`
	Name  string  `json:"name" bson:"name"`
	Lists []*List `json:"lists" bson:"lists"`
}

type Card struct {
	ID            string  `json:"_id" bson:"_id"`
	Name          string  `json:"name" bson:"name"`
	DueDate       *string `json:"dueDate" bson:"dueDate"`
	ParentBoardID string  `json:"parentBoardId" bson:"parentBoardId"`
	ParentListID  string  `json:"parentListId" bson:"parentListId"`
}

type CardIdentifier struct {
	ID            string `json:"_id" bson:"_id"`
	ParentListID  string `json:"parentListId" bson:"parentListId"`
	ParentBoardID string `json:"parentBoardId" bson:"parentBoardId"`
}

type CardMetaData struct {
	ID      string  `json:"_id" bson:"_id"`
	Name    string  `json:"name" bson:"name"`
	DueDate *string `json:"dueDate" bson:"dueDate"`
}

type ChangeCardOrder struct {
	BoardID    string `json:"boardId" bson:"boardId"`
	SrcListID  string `json:"srcListId" bson:"srcListId"`
	DestListID string `json:"destListId" bson:"destListId"`
	CardID     string `json:"cardId" bson:"cardId"`
	SrcIdx     int    `json:"srcIdx" bson:"srcIdx"`
	DestIdx    int    `json:"destIdx" bson:"destIdx"`
}

type ChangeListOrder struct {
	BoardID string `json:"boardId" bson:"boardId"`
	ListID  string `json:"listId" bson:"listId"`
	SrcIdx  int    `json:"srcIdx" bson:"srcIdx"`
	DestIdx int    `json:"destIdx" bson:"destIdx"`
}

type List struct {
	ID            string          `json:"_id" bson:"_id"`
	Name          string          `json:"name" bson:"name"`
	ParentBoardID string          `json:"parentBoardId" bson:"parentBoardId"`
	Cards         []*CardMetaData `json:"cards" bson:"cards"`
}

type ListIdentifier struct {
	ID            string `json:"_id" bson:"_id"`
	ParentBoardID string `json:"parentBoardId" bson:"parentBoardId"`
}

type NewBoard struct {
	Name string `json:"name" bson:"name"`
}

type NewCard struct {
	Name          string  `json:"name" bson:"name"`
	DueDate       *string `json:"dueDate" bson:"dueDate"`
	ParentBoardID string  `json:"parentBoardId" bson:"parentBoardId"`
	ParentListID  string  `json:"parentListId" bson:"parentListId"`
}

type NewList struct {
	Name          string `json:"name" bson:"name"`
	ParentBoardID string `json:"parentBoardId" bson:"parentBoardId"`
}
