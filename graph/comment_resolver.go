package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/olzzhas/narxozer/graph/middleware"
	"github.com/olzzhas/narxozer/graph/model"
)

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	comment := &model.Comment{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: int(userID),
		ParentID: input.ParentID,
	}

	comment, err := r.Models.Posts.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// ReplyToComment is the resolver for the replyToComment field.
func (r *mutationResolver) ReplyToComment(ctx context.Context, commentID int, input model.CreateCommentInput) (*model.Comment, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	reply := &model.Comment{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: int(userID),
		ParentID: &commentID,
		Replies:  []*model.Comment{},
	}

	reply, err := r.Models.Posts.CreateComment(reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID int) ([]*model.Comment, error) {
	comments, err := r.Models.Posts.FindAllComment(int64(postID))
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// LikeComment is the resolver for the likeComment field.
func (r *mutationResolver) LikeComment(ctx context.Context, id int) (*model.Comment, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Проверяем, не лайкнул ли уже этот пользователь данный комментарий
	var existingLike int
	err := r.Models.Posts.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'comment'", userID, id).Scan(&existingLike)
	if err != nil {
		return nil, err
	}

	if existingLike > 0 {
		return nil, fmt.Errorf("you have already liked this comment")
	}

	// Добавляем лайк в таблицу likes
	_, err = r.Models.Posts.DB.Exec("INSERT INTO likes (user_id, entity_id, entity_type) VALUES ($1, $2, 'comment')", userID, id)
	if err != nil {
		return nil, err
	}

	// Увеличиваем счетчик лайков в таблице comments
	_, err = r.Models.Posts.DB.Exec("UPDATE comments SET likes = likes + 1 WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Возвращаем обновленный комментарий
	comment := &model.Comment{}
	err = r.Models.Posts.DB.QueryRow("SELECT id, content, created_at, likes FROM comments WHERE id = $1", id).Scan(
		&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Likes)
	if err != nil {
		return nil, err
	}

	return comment, nil
	
}
