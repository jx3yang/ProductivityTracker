package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jx3yang/ProductivityTracker/src/backend/graph/generated"
	"github.com/jx3yang/ProductivityTracker/src/backend/graph/model"
)

func (r *mutationResolver) CreateCard(ctx context.Context, card model.NewCard) (*model.Card, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateList(ctx context.Context, list model.NewList) (*model.List, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCard(ctx context.Context, id string) (*model.Card, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetList(ctx context.Context, id string) (*model.List, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
