package graph

import (
	"context"
	"github.com/olzzhas/narxozer/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {

	return &model.User{}, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id int, input model.UpdateUserInput) (*model.User, error) {
	return &model.User{}, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (bool, error) {
	return true, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return []*model.User{}, nil
}

// UserByID is the resolver for the userById field.
func (r *queryResolver) UserByID(ctx context.Context, id int) (*model.User, error) {
	return &model.User{}, nil
}
