package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/jx3yang/ProductivityTracker/src/backend/graph/generated"
	"github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
	"github.com/jx3yang/ProductivityTracker/src/backend/handler"
)

func (r *mutationResolver) CreateCard(ctx context.Context, card model.NewCard) (*model.Card, error) {
	return handler.CreateCard(&card)
}

func (r *mutationResolver) CreateList(ctx context.Context, list model.NewList) (*model.List, error) {
	return handler.CreateList(&list)
}

func (r *mutationResolver) CreateBoard(ctx context.Context, board model.NewBoard) (*model.Board, error) {
	return handler.CreateBoard(&board)
}

func (r *queryResolver) GetCard(ctx context.Context, id string) (*model.Card, error) {
	return handler.FindCardByID(id)
}

func (r *queryResolver) GetList(ctx context.Context, id string) (*model.List, error) {
	return handler.FindListByID(id)
}

func (r *queryResolver) GetBoard(ctx context.Context, id string) (*model.Board, error) {
	return handler.FindBoardByID(id)
}

func (r *mutationResolver) UpdateCardOrder(ctx context.Context, changeCardOrder model.ChangeCardOrder) (bool, error) {
	return true, nil
}

func (r *mutationResolver) UpdateListOrder(ctx context.Context, changeListOrder model.ChangeListOrder) (bool, error) {
	return handler.UpdateListOrder(&changeListOrder)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
