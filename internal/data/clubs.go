package data

import (
	"database/sql"
	"fmt"
	"github.com/olzzhas/narxozer/graph/model"
)

type ClubModel struct {
	DB *sql.DB
}

func (m ClubModel) Insert(club *model.Club, id int) (*model.Club, error) {
	query := `
		INSERT INTO clubs (name, description, image_url, creator)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	args := []interface{}{club.Name, club.Description, club.ImageURL, club.Creator}

	err := m.DB.QueryRow(query, args...).Scan(&club.ID, &club.CreatedAt)
	if err != nil {
		return nil, err
	}

	return club, nil
}

func (m ClubModel) GetAll() ([]*model.Club, error) {
	query := `
		SELECT id, name, description, image_url, creator, created_at
		FROM clubs
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clubs []*model.Club
	for rows.Next() {
		var club model.Club
		err := rows.Scan(
			&club.ID,
			&club.Name,
			&club.Description,
			&club.ImageURL,
			&club.Creator,
			&club.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, &club)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clubs, nil
}

func (m ClubModel) GetByID(id int) (*model.Club, error) {
	query := `
		SELECT id, name, description, image_url, creator, created_at
		FROM clubs
		WHERE id = $1`

	var club model.Club
	err := m.DB.QueryRow(query, id).Scan(
		&club.ID,
		&club.Name,
		&club.Description,
		&club.ImageURL,
		&club.Creator,
		&club.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пост не найден
		}
		return nil, err
	}

	return &club, nil
}

func (m ClubModel) AddMember(clubID, userID int) error {
	query := `
		INSERT INTO club_members (club_id, user_id)
		VALUES ($1, $2)
	`

	err := m.DB.QueryRow(query, clubID, userID)
	if err != nil {
		return fmt.Errorf("error occured while writing data into db: %v", err)
	}

	return nil
}

func (m ClubModel) RemoveMember(clubID int, userID int) error {
	query := `
		DELETE FROM club_members 
		WHERE club_id = $1 and user_id = $2
	`

	err := m.DB.QueryRow(query, clubID, userID)
	if err != nil {
		return fmt.Errorf("error occured while deleting data from db: %v", err)
	}

	return nil
}

func (m ClubModel) IsAdmin(clubId, userId int) bool {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM club_admins WHERE club_id = $1 AND user_id = $2
        )`

	err := m.DB.QueryRow(query, clubId, userId).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}
