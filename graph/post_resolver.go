package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/olzzhas/narxozer/graph/middleware"
	"github.com/olzzhas/narxozer/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// PostByID is the resolver for the postById field.
func (r *queryResolver) PostByID(ctx context.Context, id int) (*model.Post, error) {

	post, err := r.Models.Posts.FindOne(int64(id))
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	if post == nil {
		return nil, gqlerror.Errorf("post not found")
	}

	// TODO redis
	user, err := r.Models.Users.Get(post.Author.ID)
	if err != nil {
		return nil, gqlerror.Errorf("internal server error")
	}

	post.Author = user

	return post, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.Models.Posts.FindAll()
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while getting posts: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	for _, post := range posts {
		//TODO redis
		user, err := r.Models.Users.Get(post.Author.ID)
		if err != nil {
			r.Logger.PrintError(fmt.Errorf("error while getting user: %v", err), nil)
			return nil, gqlerror.Errorf("internal server error")
		}
		post.Author = user
	}

	return posts, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	temp := model.Post{
		Title:    input.Title,
		Content:  input.Content,
		ImageURL: input.ImageURL,
		Author:   &model.User{ID: int(userID)},
	}

	post, err := r.Models.Posts.Insert(&temp)
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	//TODO redis
	user, err := r.Models.Users.Get(post.Author.ID)
	if err != nil {
		return nil, gqlerror.Errorf("internal server error")
	}

	post.Author = user

	return post, nil
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, id int, input model.UpdatePostInput) (*model.Post, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	// Получаем пост, чтобы обновить его поля
	post, err := r.Models.Posts.FindOne(int64(id))
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	if post == nil {
		return nil, gqlerror.Errorf("post not found")
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.ImageURL != nil {
		post.ImageURL = input.ImageURL
	}

	err = r.Models.Posts.Update(post)
	if err != nil {
		r.Logger.PrintError(err, nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	//TODO redis
	user, err := r.Models.Users.Get(post.Author.ID)
	if err != nil {
		return nil, gqlerror.Errorf("internal server error")
	}

	post.Author = user

	return post, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, id int) (bool, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return false, errors.New("unauthorized")
	}

	// TODO redis
	post, err := r.Models.Posts.FindOne(int64(id))
	if err != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	if post == nil {
		return false, gqlerror.Errorf("post not found")
	}

	//TODO redis
	user, err := r.Models.Users.Get(int(userID))
	if err != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	if post.Author.ID != int(userID) && user.Role != model.RoleAdmin {
		return false, gqlerror.Errorf("you have no permission to delete this post")
	}

	err = r.Models.Posts.Delete(int64(id))
	if err != nil {
		r.Logger.PrintError(err, nil)
		return false, gqlerror.Errorf("internal server error")
	}

	return true, nil
}

// LikePost is the resolver for the likePost field.
func (r *mutationResolver) LikePost(ctx context.Context, id int) (*model.Post, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Проверяем, не лайкнул ли уже этот пользователь данный пост
	var existingLike int
	err := r.Models.Posts.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'post'", userID, id).Scan(&existingLike)
	if err != nil {
		return nil, err
	}

	if existingLike > 0 {
		// Если пользователь уже лайкнул пост, удаляем лайк (unlike)
		_, err := r.Models.Posts.DB.Exec("DELETE FROM likes WHERE user_id = $1 AND entity_id = $2 AND entity_type = 'post'", userID, id)
		if err != nil {
			return nil, gqlerror.Errorf("internal server error")
		}

		// Уменьшаем счетчик лайков в таблице posts
		_, err = r.Models.Posts.DB.Exec("UPDATE posts SET likes = likes - 1 WHERE id = $1", id)
		if err != nil {
			return nil, gqlerror.Errorf("internal server error")
		}

		// Возвращаем обновленный пост
		post := &model.Post{}
		err = r.Models.Posts.DB.QueryRow("SELECT id, title, content, image_url, created_at, updated_at, likes FROM posts WHERE id = $1", id).Scan(
			&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt, &post.Likes)
		if err != nil {
			return nil, gqlerror.Errorf("internal server error")
		}

		return post, nil
	}

	// Добавляем лайк в таблицу likes
	_, err = r.Models.Posts.DB.Exec("INSERT INTO likes (user_id, entity_id, entity_type) VALUES ($1, $2, 'post')", userID, id)
	if err != nil {
		return nil, err
	}

	// Увеличиваем счетчик лайков в таблице posts
	_, err = r.Models.Posts.DB.Exec("UPDATE posts SET likes = likes + 1 WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Возвращаем обновленный пост
	post := &model.Post{}
	err = r.Models.Posts.DB.QueryRow("SELECT id, title, content, image_url, created_at, updated_at, likes FROM posts WHERE id = $1", id).Scan(
		&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt, &post.Likes)
	if err != nil {
		return nil, err
	}

	return post, nil

}
