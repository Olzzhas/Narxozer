package data

import (
	"database/sql"
	"github.com/olzzhas/narxozer/graph/model"
)

type CommentModel struct {
	DB *sql.DB
}

func (m CommentModel) Insert(comment *model.Comment) (*model.Comment, error) {
	query := `
		INSERT INTO comments (content, image_url, entity_id, entity_type, author_id, parent_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
		RETURNING id, created_at
	`

	args := []interface{}{comment.Content, comment.ImageURL, comment.EntityID, comment.EntityType, comment.ParentID}

	err := m.DB.QueryRow(query, args...).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (m CommentModel) GetAllByPost(id int) ([]*model.Comment, error) {
	query := `
		SELECT id, content, image_url, entity_id, entity_type, author_id, parent_id, created_at, updated_at, likes
		FROM comments
		WHERE entity_id = $1
	`

	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		comment.Author = &model.User{}
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.ImageURL,
			&comment.EntityID,
			&comment.EntityType,
			&comment.Author.ID,
			&comment.ParentID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Likes,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (m CommentModel) GetByEntityID(topicID int, entityType string) ([]*model.Comment, error) {
	query := `
		SELECT id, content, entity_id, entity_type, author_id, parent_id, created_at, updated_at, likes
		FROM comments
		WHERE entity_id = $1 AND entity_type = $2
		ORDER BY created_at ASC`

	rows, err := m.DB.Query(query, topicID, entityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.EntityID,
			&comment.EntityType,
			&comment.Author.ID,
			&comment.ParentID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Likes,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (m CommentModel) UpdatePostComment() {

}

func (m CommentModel) DeletePostComment() {

}

func (m CommentModel) GetPostComments() {

}

func (m CommentModel) InsertTopicComment() {

}

func (m CommentModel) UpdateTopicComment() {

}

func (m CommentModel) DeleteTopicComment() {

}

func (m CommentModel) GetTopicComments() {

}
