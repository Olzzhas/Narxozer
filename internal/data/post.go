package data

import (
	"database/sql"
	"github.com/olzzhas/narxozer/graph/model"
)

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) Insert(post *model.CreatePostInput) error {
	query := `
		INSERT INTO lessons (title, content, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, version
	`
}
