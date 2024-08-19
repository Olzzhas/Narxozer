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

	args := []interface{}{club.Name, club.Description, club.ImageURL, club.Creator.ID}

	err := m.DB.QueryRow(query, args...).Scan(&club.ID, &club.CreatedAt)
	if err != nil {
		return nil, err
	}

	return club, nil
}

func (m ClubModel) GetAll() ([]*model.Club, error) {
	query := `
		SELECT id, name, description, image_url, creator_id, created_at
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
		club.Creator = &model.User{}
		err := rows.Scan(
			&club.ID,
			&club.Name,
			&club.Description,
			&club.ImageURL,
			&club.Creator.ID,
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
		SELECT id, name, description, image_url, creator_id, created_at
		FROM clubs
		WHERE id = $1`

	var club model.Club
	club.Creator = &model.User{}
	err := m.DB.QueryRow(query, id).Scan(
		&club.ID,
		&club.Name,
		&club.Description,
		&club.ImageURL,
		&club.Creator.ID,
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

func (m ClubModel) IsCreator(clubId, userId int) (bool, error) {
	query := `SELECT COUNT(*) FROM clubs WHERE id = $1 AND creator_id = $2`
	var count int
	err := m.DB.QueryRow(query, clubId, userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m ClubModel) Delete(id int) error {
	query := `DELETE FROM clubs WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m ClubModel) DeleteAllRelatedData(id int) error {
	_, err := m.DB.Exec(`DELETE FROM events WHERE club_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(`DELETE FROM club_admins WHERE club_id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (m ClubModel) AddAdmin(clubID, userID int) error {
	query := `INSERT INTO club_admins (club_id, user_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(query, clubID, userID)
	return err
}
