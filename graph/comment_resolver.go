package graph

import (
	"context"
	"github.com/olzzhas/narxozer/graph/model"
)

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	comment := &model.Comment{
		ID:        "some-generated-id", // Логика для генерации ID
		Content:   input.Content,
		CreatedAt: "some-timestamp", // Логика для установки времени
		Likes:     0,
		Replies:   []*model.Comment{},
	}

	// Логика для сохранения комментария
	// Здесь должна быть реализация сохранения комментария в базу данных

	return comment, nil
}

// LikeComment is the resolver for the likeComment field.
func (r *mutationResolver) LikeComment(ctx context.Context, id string) (*model.Comment, error) {
	//userID := 1 // Пример получения userID, вам нужно получить его из контекста или токена
	//
	//// Проверяем, не лайкнул ли уже этот пользователь данный комментарий
	//var existingLike int
	//err := r.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'comment'", userID, id).Scan(&existingLike)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if existingLike > 0 {
	//	return nil, fmt.Errorf("you have already liked this comment")
	//}
	//
	//// Добавляем лайк в таблицу likes
	//_, err = r.DB.Exec("INSERT INTO likes (user_id, entity_id, entity_type) VALUES ($1, $2, 'comment')", userID, id)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Увеличиваем счетчик лайков в таблице comments
	//_, err = r.DB.Exec("UPDATE comments SET likes = likes + 1 WHERE id = $1", id)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Возвращаем обновленный комментарий
	//comment := &model.Comment{}
	//err = r.DB.QueryRow("SELECT id, content, created_at, likes FROM comments WHERE id = $1", id).Scan(
	//	&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Likes)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return comment, nil

	return nil, nil
}

// ReplyToComment is the resolver for the replyToComment field.
func (r *mutationResolver) ReplyToComment(ctx context.Context, commentID string, input model.CreateCommentInput) (*model.Comment, error) {
	reply := &model.Comment{
		ID:        "some-generated-id", // Логика для генерации ID
		Content:   input.Content,
		CreatedAt: "some-timestamp", // Логика для установки времени
		Likes:     0,
		Replies:   []*model.Comment{},
	}

	// Логика для сохранения ответа на комментарий
	// Здесь должна быть реализация сохранения ответа в базу данных

	return reply, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string) ([]*model.Comment, error) {
	return []*model.Comment{}, nil
}
