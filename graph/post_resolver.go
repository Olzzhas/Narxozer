package graph

import (
	"context"
	"github.com/olzzhas/narxozer/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"strconv"
)

// PostByID is the resolver for the postById field.
func (r *queryResolver) PostByID(ctx context.Context, id string) (*model.Post, error) {
	postID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gqlerror.Errorf("invalid post ID")
	}

	post, err := r.Models.Posts.FindOne(postID)
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	if post == nil {
		return nil, gqlerror.Errorf("post not found")
	}

	return post, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.Models.Posts.FindAll()
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	return posts, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	post, err := r.Models.Posts.Insert(&input)
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	return post, nil
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input model.UpdatePostInput) (*model.Post, error) {
	postID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gqlerror.Errorf("invalid post ID")
	}

	// Получаем пост, чтобы обновить его поля
	post, err := r.Models.Posts.FindOne(postID)
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	if post == nil {
		return nil, gqlerror.Errorf("post not found")
	}

	// Обновляем поля поста
	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.ImageURL != nil {
		post.ImageURL = input.ImageURL
	}

	// Сохраняем обновления
	err = r.Models.Posts.Update(post)
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	return post, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	postID, err := strconv.ParseInt(id, 10, 64) // Преобразуем строку в int64
	if err != nil {
		return false, gqlerror.Errorf("invalid post ID")
	}

	err = r.Models.Posts.Delete(postID)
	if err != nil {
		r.Logger.PrintError(err, nil)
		return false, gqlerror.Errorf("internal server error")
	}

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
