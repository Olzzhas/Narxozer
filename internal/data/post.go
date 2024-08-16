package data

import (
	"database/sql"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert()
