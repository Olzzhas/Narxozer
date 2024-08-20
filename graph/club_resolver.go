package graph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olzzhas/narxozer/graph/middleware"
	"github.com/olzzhas/narxozer/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"time"
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

	// Получаем обновленный список членов клуба
	members, err := r.Models.Clubs.GetMembers(clubID)
	if err != nil {
		return nil, err
	}

	// Обновляем кеш в Redis
	cacheKey := fmt.Sprintf("club:%d:members", clubID)
	data, err := json.Marshal(members)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while marshaling club members: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	err = r.Models.Clubs.Redis.Set(ctx, cacheKey, data, 10*time.Minute).Err()
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while updating cache: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	// Возвращаем обновленные данные о клубе
	club, err := r.Models.Clubs.GetCachedByID(clubID)
	if err != nil {
		return nil, err
	}
	club.Members = members

	return club, nil
}

func (r *mutationResolver) LeaveClub(ctx context.Context, clubID int) (*model.Club, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	success, err := r.Models.Clubs.RemoveMember(clubID, int(userID))
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, errors.New("failed to remove member from club")
	}

	// Получаем обновленный список членов клуба
	members, err := r.Models.Clubs.GetMembers(clubID)
	if err != nil {
		return nil, err
	}

	// Обновляем кеш в Redis
	cacheKey := fmt.Sprintf("club:%d:members", clubID)
	data, err := json.Marshal(members)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while marshaling club members: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	err = r.Models.Clubs.Redis.Set(ctx, cacheKey, data, 10*time.Minute).Err()
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while updating cache: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	// Возвращаем обновленные данные о клубе
	club, err := r.Models.Clubs.GetCachedByID(clubID)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while getting club: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}
	club.Members = members

	return club, nil
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
		Creator:     &model.User{ID: int(userID)},
	}

	newClub, err := r.Models.Clubs.Insert(club, int(userID))
	if err != nil {
		return nil, err
	}

	err = r.Models.Clubs.AddAdmin(newClub.ID, int(userID))
	if err != nil {
		return nil, err
	}

	user, err := r.Models.Users.GetCached(int(userID))
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while getting user: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	newClub.Creator = user

	return newClub, nil
}

// UpdateClub is the resolver for the updateClub field.
func (r *mutationResolver) UpdateClub(ctx context.Context, id int, input model.UpdateClubInput) (*model.Club, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	isAdmin := r.Models.Clubs.IsAdmin(id, int(userID))
	if !isAdmin {
		return nil, errors.New("you do not have permission to update this club")
	}

	club, err := r.Models.Clubs.Update(id, input)
	if err != nil {
		return nil, err
	}

	user, err := r.Models.Users.GetCached(int(userID))
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while getting user: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	club.Creator = user

	// Обновляем данные в кеше Redis
	cacheKey := fmt.Sprintf("club:%d", id)
	data, err := json.Marshal(club)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while marshaling club data: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	err = r.Models.Clubs.Redis.Set(ctx, cacheKey, data, 10*time.Minute).Err()
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while updating cache: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	return club, nil
}

// DeleteClub is the resolver for the deleteClub field.
func (r *mutationResolver) DeleteClub(ctx context.Context, id int) (bool, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return false, errors.New("unauthorized")
	}

	isAdmin := r.Models.Clubs.IsAdmin(id, int(userID))

	// Если пользователь не является администратором, проверяем, является ли он создателем клуба
	if !isAdmin {
		isCreator, err := r.Models.Clubs.IsCreator(id, int(userID))
		if err != nil {
			return false, err
		}
		if !isCreator {
			return false, errors.New("you do not have permission to delete this club")
		}
	}

	// Удаляем все связанные с клубом данные
	err := r.Models.Clubs.DeleteAllRelatedData(id)
	if err != nil {
		return false, err
	}

	// Удаляем сам клуб
	err = r.Models.Clubs.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Clubs is the resolver for the clubs field.
func (r *queryResolver) Clubs(ctx context.Context) ([]*model.Club, error) {
	clubs, err := r.Models.Clubs.GetAll()
	if err != nil {
		return nil, err
	}
	for _, club := range clubs {
		members, err := r.Models.Clubs.GetCachedMembers(club.ID)
		if err != nil {
			r.Logger.PrintError(fmt.Errorf("error while getting club members: %v", err), nil)
			return nil, gqlerror.Errorf("internal server error")
		}
		club.Members = members
	}

	return clubs, nil
}

// ClubByID is the resolver for the clubById field.
func (r *queryResolver) ClubByID(ctx context.Context, id int) (*model.Club, error) {
	club, err := r.Models.Clubs.GetCachedByID(id)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while getting club: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	members, err := r.Models.Clubs.GetCachedMembers(club.ID)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while getting club members: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}
	club.Members = members

	return club, nil

}

// AssignAdmin is the resolver for the assignAdmin field.
func (r *mutationResolver) AssignAdmin(ctx context.Context, clubID int, userID int) (*model.Club, error) {
	adminID := middleware.GetUserIDFromContext(ctx)
	if adminID == 0 {
		return nil, errors.New("unauthorized")
	}

	isAdmin := r.Models.Clubs.IsAdmin(clubID, int(adminID))
	if !isAdmin {
		return nil, errors.New("you do not have permission to assign admins for this club")
	}

	err := r.Models.Clubs.AddAdmin(clubID, userID)
	if err != nil {
		return nil, err
	}

	club, err := r.Models.Clubs.GetByID(clubID)
	if err != nil {
		return nil, err
	}

	user, err := r.Models.Users.GetCached(club.Creator.ID)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while getting user: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	club.Creator = user

	return club, nil
}
