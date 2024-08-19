package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/olzzhas/narxozer/graph/middleware"
	"github.com/olzzhas/narxozer/graph/model"
)

// CreateTopic is the resolver for the createTopic field.
func (r *mutationResolver) CreateTopic(ctx context.Context, input model.CreateTopicInput) (*model.Topic, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	topic := &model.Topic{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: int(userID),
	}

	topic, err := r.Models.Topics.Insert(topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// UpdateTopic is the resolver for the updateTopic field.
func (r *mutationResolver) UpdateTopic(ctx context.Context, id int, input model.UpdateTopicInput) (*model.Topic, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	topic, err := r.Models.Topics.GetByID(id)
	if err != nil {
		return nil, err
	}

	if topic.AuthorID != int(userID) {
		return nil, errors.New("unauthorized")
	}

	if input.Title != nil {
		topic.Title = *input.Title
	}
	if input.Content != nil {
		topic.Content = *input.Content
	}

	topic, err = r.Models.Topics.Update(topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// DeleteTopic is the resolver for the deleteTopic field.
func (r *mutationResolver) DeleteTopic(ctx context.Context, id int) (bool, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return false, errors.New("unauthorized")
	}

	topic, err := r.Models.Topics.GetByID(id)
	if err != nil {
		return false, err
	}

	if topic.AuthorID != int(userID) {
		return false, errors.New("unauthorized")
	}

	err = r.Models.Topics.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// LikeTopic is the resolver for the likeTopic field.
func (r *mutationResolver) LikeTopic(ctx context.Context, id int) (*model.Topic, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Проверяем, не лайкнул ли уже этот пользователь данный пост
	var existingLike int
	err := r.Models.Topics.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'topic'", userID, id).Scan(&existingLike)
	if err != nil {
		return nil, err
	}

	if existingLike > 0 {
		return nil, fmt.Errorf("you have already liked this post")
	}

	// Добавляем лайк в таблицу likes
	_, err = r.Models.Topics.DB.Exec("INSERT INTO likes (user_id, entity_id, entity_type) VALUES ($1, $2, 'topic')", userID, id)
	if err != nil {
		return nil, err
	}

	// Увеличиваем счетчик лайков в таблице posts
	_, err = r.Models.Topics.DB.Exec("UPDATE topics SET likes = likes + 1 WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Возвращаем обновленный пост
	topic := &model.Topic{}
	err = r.Models.Posts.DB.QueryRow("SELECT id, title, content, author_id, created_at, updated_at, likes FROM posts WHERE id = $1", id).Scan(
		&topic.ID, &topic.Title, &topic.Content, &topic.AuthorID, &topic.CreatedAt, &topic.UpdatedAt, &topic.Likes)
	if err != nil {
		return nil, err
	}

	return topic, nil

}

// UpdateComment is the resolver for the updateComment field.
func (r *mutationResolver) UpdateComment(ctx context.Context, id int, input model.UpdateCommentInput) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented: UpdateComment - updateComment"))
}

// DeleteComment is the resolver for the deleteComment field.
func (r *mutationResolver) DeleteComment(ctx context.Context, id int) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteComment - deleteComment"))
}

// Topics is the resolver for the topics field.
func (r *queryResolver) Topics(ctx context.Context) ([]*model.Topic, error) {
	topics, err := r.Models.Topics.GetAll()
	if err != nil {
		return nil, err
	}
	return topics, nil
}

// TopicByID is the resolver for the topicById field.
func (r *queryResolver) TopicByID(ctx context.Context, id int) (*model.Topic, error) {
	topic, err := r.Models.Topics.GetByID(id)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

// CommentsByTopicID is the resolver for the commentsByTopicId field.
func (r *queryResolver) CommentsByTopicID(ctx context.Context, topicID int) ([]*model.Comment, error) {
	// TODO implement

	panic(fmt.Errorf("not implemented: CommentsByTopicID - commentsByTopicId"))
}
