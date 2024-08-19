package graph

import (
	"context"
	"errors"
	"github.com/olzzhas/narxozer/graph/middleware"
	"github.com/olzzhas/narxozer/graph/model"
)

// CreateEvent is the resolver for the createEvent field.
func (r *mutationResolver) CreateEvent(ctx context.Context, clubID int, input model.CreateEventInput) (*model.Event, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	if !r.Models.Clubs.IsAdmin(clubID, int(userID)) {
		return nil, errors.New("unauthorized: only admins can create events")
	}

	event := &model.Event{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		ClubID:      clubID,
	}

	event, err := r.Models.Events.Insert(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// UpdateEvent is the resolver for the updateEvent field.
func (r *mutationResolver) UpdateEvent(ctx context.Context, id int, input model.UpdateEventInput) (*model.Event, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	event, err := r.Models.Events.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Проверяем, является ли пользователь администратором клуба
	if !r.Models.Clubs.IsAdmin(event.ClubID, int(userID)) {
		return nil, errors.New("unauthorized: only admins can update events")
	}

	// Обновляем поля мероприятия, если они были переданы в input
	if input.Title != nil {
		event.Title = *input.Title
	}
	if input.Description != nil {
		event.Description = *input.Description
	}
	if input.Date != nil {
		event.Date = *input.Date
	}

	event, err = r.Models.Events.Update(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// DeleteEvent is the resolver for the deleteEvent field.
func (r *mutationResolver) DeleteEvent(ctx context.Context, id int) (bool, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return false, errors.New("unauthorized")
	}

	event, err := r.Models.Events.GetByID(id)
	if err != nil {
		return false, err
	}

	// Проверяем, является ли пользователь администратором клуба
	if !r.Models.Clubs.IsAdmin(event.ClubID, int(userID)) {
		return false, errors.New("unauthorized: only admins can delete events")
	}

	err = r.Models.Events.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
