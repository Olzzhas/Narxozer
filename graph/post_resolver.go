package graph

import (
	"context"
	"github.com/olzzhas/narxozer/graph/model"
)

// PostByID is the resolver for the postById field.
func (r *queryResolver) PostByID(ctx context.Context, id string) (*model.Post, error) {
	return &model.Post{}, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {

	return []*model.Post{}, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	post := &model.Post{
		ID:        "some-generated-id", // Логика для генерации ID
		Title:     input.Title,
		Content:   input.Content,
		AuthorID:  "some-generated-id",
		ImageURL:  input.ImageURL,   // ImageURL может быть nil
		CreatedAt: "some-timestamp", // Логика для установки времени
		Likes:     0,
		Comments:  []*model.Comment{},
	}

	// Логика для сохранения поста
	// Здесь должна быть реализация сохранения поста в базу данных

	return post, nil
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input model.UpdatePostInput) (*model.Post, error) {
	// Логика для обновления поста
	// Здесь должна быть реализация обновления поста в базе данных
	return &model.Post{}, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	// Логика для удаления поста
	// Здесь должна быть реализация удаления поста из базы данных
	return true, nil
}

// LikePost is the resolver for the likePost field.
func (r *mutationResolver) LikePost(ctx context.Context, id string) (*model.Post, error) {
	//userID := 1 // Пример получения userID, вам нужно получить его из контекста или токена
	//
	//// Проверяем, не лайкнул ли уже этот пользователь данный пост
	//var existingLike int
	//err := r.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'post'", userID, id).Scan(&existingLike)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if existingLike > 0 {
	//	return nil, fmt.Errorf("you have already liked this post")
	//}
	//
	//// Добавляем лайк в таблицу likes
	//_, err = r.DB.Exec("INSERT INTO likes (user_id, entity_id, entity_type) VALUES ($1, $2, 'post')", userID, id)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Увеличиваем счетчик лайков в таблице posts
	//_, err = r.DB.Exec("UPDATE posts SET likes = likes + 1 WHERE id = $1", id)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Возвращаем обновленный пост
	//post := &model.Post{}
	//err = r.DB.QueryRow("SELECT id, title, content, image_url, created_at, updated_at, likes FROM posts WHERE id = $1", id).Scan(
	//	&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt, &post.Likes)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return post, nil

	return nil, nil
}
