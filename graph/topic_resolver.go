package graph

import (
	"context"
	"errors"
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
		ImageURL: input.ImageURL,
		Author:   &model.User{ID: int(userID)},
	}

	topic, err := r.Models.Topics.Insert(topic)
	if err != nil {
		return nil, err
	}

	// TODO redis
	user, err := r.Models.Users.Get(int(userID))
	if err != nil {
		return nil, err
	}

	topic.Author = user

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
	if topic == nil {
		return nil, errors.New("topic not found")
	}

	if topic.Author.ID != int(userID) {
		return nil, errors.New("you do not have permission to update this topic")
	}

	topic.Title = *input.Title
	topic.Content = *input.Content
	topic.ImageURL = input.ImageURL

	updatedTopic, err := r.Models.Topics.Update(topic)
	if err != nil {
		return nil, err
	}

	// TODO redis
	user, err := r.Models.Users.Get(int(userID))
	if err != nil {
		return nil, err
	}

	updatedTopic.Author = user

	return updatedTopic, nil
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
	if topic == nil {
		return false, errors.New("topic not found")
	}

	// Проверяем, является ли пользователь автором топика
	if topic.Author.ID != int(userID) {
		return false, errors.New("you do not have permission to delete this topic")
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
		return nil, errors.New("unauthorized")
	}

	// Проверяем, лайкнул ли уже этот пользователь данный топик
	var existingLike int
	err := r.Models.Posts.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'topic'", userID, id).Scan(&existingLike)
	if err != nil {
		return nil, err
	}

	tx, err := r.Models.Posts.DB.Begin() // Начало транзакции
	if err != nil {
		return nil, err
	}

	if existingLike > 0 {
		// Лайк уже существует, выполняем анлайк
		_, err = tx.Exec("DELETE FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'topic'", userID, id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Уменьшаем счетчик лайков в таблице topics, убедившись, что он не станет отрицательным
		_, err = tx.Exec("UPDATE topics SET likes = GREATEST(likes - 1, 0) WHERE id = $1", id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		// Лайка нет, выполняем добавление лайка
		_, err = tx.Exec("INSERT INTO likes (user_id, entity_id, entity_type) VALUES ($1, $2, 'topic')", userID, id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Увеличиваем счетчик лайков в таблице topics
		_, err = tx.Exec("UPDATE topics SET likes = likes + 1 WHERE id = $1", id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil { // Завершение транзакции
		return nil, err
	}

	// Возвращаем обновленный топик
	topic, err := r.Models.Topics.GetByID(id)
	if err != nil {
		return nil, err
	}

	// TODO redis
	user, err := r.Models.Users.Get(int(userID))
	if err != nil {
		return nil, err
	}

	topic.Author = user

	return topic, nil
}

// Topics is the resolver for the topics field.
func (r *queryResolver) Topics(ctx context.Context) ([]*model.Topic, error) {
	topics, err := r.Models.Topics.GetAll()
	if err != nil {
		return nil, err
	}

	// Загружаем авторов для каждого топика
	for _, topic := range topics {
		user, err := r.Models.Users.Get(topic.Author.ID)
		if err != nil {
			return nil, err
		}
		topic.Author = user
	}

	return topics, nil
}

// TopicByID is the resolver for the topicById field.
func (r *queryResolver) TopicByID(ctx context.Context, id int) (*model.Topic, error) {
	topic, err := r.Models.Topics.GetByID(id)
	if err != nil {
		return nil, err
	}
	if topic == nil {
		return nil, errors.New("topic not found")
	}

	// Загружаем автора топика
	user, err := r.Models.Users.Get(topic.Author.ID)
	if err != nil {
		return nil, err
	}
	topic.Author = user

	return topic, nil
}

// CommentsByTopicID is the resolver for the commentsByTopicId field.
func (r *queryResolver) CommentsByTopicID(ctx context.Context, topicID int) ([]*model.Comment, error) {
	comments, err := r.Models.Comments.GetByEntityID(topicID, "topic")
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		// TODO redis
		user, err := r.Models.Users.Get(comment.Author.ID)
		if err != nil {
			return nil, err
		}
		comment.Author = user
	}

	return comments, nil
}
