package graph

import (
	"context"
	"errors"
	"github.com/olzzhas/narxozer/graph/middleware"
	"github.com/olzzhas/narxozer/graph/model"
)

// JoinClub is the resolver for the joinClub field.
func (r *mutationResolver) JoinClub(ctx context.Context, clubID int) (*model.Club, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	err := r.Models.Clubs.AddMember(clubID, int(userID))
	if err != nil {
		return nil, err
	}

	return r.Models.Clubs.GetByID(clubID)
}

// LeaveClub is the resolver for the leaveClub field.
func (r *mutationResolver) LeaveClub(ctx context.Context, clubID int) (*model.Club, error) {
	// TODO implement
	return nil, nil
}

// CreateClub is the resolver for the createClub field.
func (r *mutationResolver) CreateClub(ctx context.Context, input model.CreateClubInput) (*model.Club, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	club := &model.Club{
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
	}

	newClub, err := r.Models.Clubs.Insert(club, int(userID))
	if err != nil {
		return nil, err
	}

	return newClub, nil
}

// UpdateClub is the resolver for the updateClub field.
func (r *mutationResolver) UpdateClub(ctx context.Context, id int, input model.UpdateClubInput) (*model.Club, error) {
	// TODO implement
	return nil, nil
}

// DeleteClub is the resolver for the deleteClub field.
func (r *mutationResolver) DeleteClub(ctx context.Context, id int) (bool, error) {
	// TODO implement
	return false, nil
}

// Clubs is the resolver for the clubs field.
func (r *queryResolver) Clubs(ctx context.Context) ([]*model.Club, error) {
	clubs, err := r.Models.Clubs.GetAll()
	if err != nil {
		return nil, err
	}
	return clubs, nil
}

// ClubByID is the resolver for the clubById field.
func (r *queryResolver) ClubByID(ctx context.Context, id int) (*model.Club, error) {
	club, err := r.Models.Clubs.GetByID(id)
	if err != nil {
		return nil, err
	}
	return club, nil
}
