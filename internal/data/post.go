package data

import (
	"database/sql"
	"github.com/olzzhas/narxozer/graph/model"
)

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) Insert(createPost *model.CreatePostInput) (*model.Post, error) {
	post := &model.Post{
		Title:     createPost.Title,
		Content:   createPost.Content,
		ImageURL:  createPost.ImageURL,
		AuthorID:  createPost.AuthorID,
		CreatedAt: "some-timestamp",
	}

	query := `
		INSERT INTO posts (title, content, image_url, author_id, created_at)
		VALUES ($1, $2, $3, $4, now())
		RETURNING id, created_at`

	args := []interface{}{createPost.Title, createPost.Content, createPost.ImageURL, createPost.AuthorID}

	err := m.DB.QueryRow(query, args...).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (m PostModel) FindOne(id int64) (*model.Post, error) {
	query := `
		SELECT id, title, content, image_url, author_id, created_at, updated_at, likes
		FROM posts
		WHERE id = $1`

	var post model.Post
	err := m.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.ImageURL,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Likes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пост не найден
		}
		return nil, err
	}

	return &post, nil
}

func (m PostModel) FindAll() ([]*model.Post, error) {
	query := `
		SELECT id, title, content, image_url, author_id, created_at, updated_at, likes
		FROM posts`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.ImageURL,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Likes,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m PostModel) Update(post *model.Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, image_url = $3, updated_at = now()
		WHERE id = $4`

	args := []interface{}{post.Title, post.Content, post.ImageURL, post.ID}

	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m PostModel) Delete(id int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m PostModel) CreateComment(comment *model.CreateCommentInput) (*model.Comment, error) {
	query := `
		INSERT INTO comments (post_id, content, author_id, created_at)
		VALUES ($1, $2, $3, now())
		RETURNING id, created_at`

	var newComment model.Comment
	args := []interface{}{comment.PostID, comment.Content, comment.AuthorID}

	err := m.DB.QueryRow(query, args...).Scan(&newComment.ID, &newComment.CreatedAt)
	if err != nil {
		return nil, err
	}

	newComment.Content = comment.Content
	newComment.AuthorID = comment.AuthorID

	return &newComment, nil
}

func (m PostModel) DeleteComment(id int64) error {
	query := `
		DELETE FROM comments
		WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m PostModel) UpdateComment(comment *model.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, updated_at = now()
		WHERE id = $2`

	args := []interface{}{comment.Content, comment.ID}

	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m PostModel) FindAllComment(postID int64) ([]*model.Comment, error) {
	query := `
		SELECT id, post_id, content, author_id, created_at, updated_at, likes
		FROM comments
		WHERE post_id = $1`

	rows, err := m.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.AuthorID,
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

func (m PostModel) FindOneComment(id int64) (*model.Comment, error) {
	query := `
		SELECT id, post_id, content, author_id, created_at, updated_at, likes
		FROM comments
		WHERE id = $1`

	var comment model.Comment
	err := m.DB.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.Content,
		&comment.AuthorID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Likes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Комментарий не найден
		}
		return nil, err
	}

	return &comment, nil
}
