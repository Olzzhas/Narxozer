package data

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/olzzhas/narxozer/graph/model"
)

type EventModel struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (m EventModel) Insert(event *model.Event) (*model.Event, error) {
	query := `
		INSERT INTO events (title, description, date, club_id, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	args := []interface{}{event.Title, event.Description, event.Date, event.ClubID, event.ImageURL}

	err := m.DB.QueryRow(query, args...).Scan(&event.ID, &event.CreatedAt)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (m EventModel) GetByID(id int) (*model.Event, error) {
	query := `
		SELECT id, title, description, date, club_id, created_at
		FROM events
		WHERE id = $1`

	var event model.Event
	err := m.DB.QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.ClubID,
		&event.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пост не найден
		}
		return nil, err
	}

	return &event, nil
}

func (m EventModel) Update(event *model.Event) (*model.Event, error) {
	query := `
		UPDATE events
		SET title = $1, description = $2, date = $3, image_url = $4
		WHERE id = $5`

	args := []interface{}{event.Title, event.Description, event.Date, event.ImageURL, event.ID}

	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return event, err
}

func (m EventModel) Delete(id int) error {
	query := `
		DELETE FROM events
		WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
