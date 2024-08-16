package data

import (
	"database/sql"
	"time"
)

type PostModel struct {
	DB *sql.DB
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	Author    User      `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID              int       `json:"id"`
	NewsID          int       `json:"news_id"`
	Content         string    `json:"content"`
	AuthorID        int       `json:"author_id"`
	Author          User      `json:"author"`
	ParentCommentID *int      `json:"parent_comment_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type Like struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
