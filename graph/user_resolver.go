package graph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/olzzhas/narxozer/graph/middleware"
	"github.com/olzzhas/narxozer/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"time"
)

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id int, input model.UpdateUserInput) (*model.User, error) {
	userId := middleware.GetUserIDFromContext(ctx)
	if userId == 0 || id != int(userId) {
		return nil, gqlerror.Errorf("you have no permission to update this user")
	}

	// Обновляем данные пользователя в базе данных
	user, err := r.Models.Users.Update(id, input)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while updating user: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	// Обновляем данные в кеше Redis
	cacheKey := fmt.Sprintf("user:%d", id)
	data, err := json.Marshal(user)
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while marshaling user data: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	err = r.Models.Users.Redis.Set(ctx, cacheKey, data, 10*time.Minute).Err()
	if err != nil {
		r.Logger.PrintError(fmt.Errorf("error while updating cache: %v", err), nil)
		return nil, gqlerror.Errorf("internal server error")
	}

	return user, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	cacheKey := "users:all"

	// Пытаемся получить данные из кеша Redis
	val, err := r.Models.Users.Redis.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		// Данные не найдены в кеше, загружаем из базы данных
		users, err := r.Models.Users.GetAll()
		if err != nil {
			return nil, err
		}

		// Сохраняем данные в кеш Redis
		data, err := json.Marshal(users)
		if err != nil {
			return nil, err
		}
		err = r.Models.Users.Redis.Set(ctx, cacheKey, data, 10*time.Minute).Err()
		if err != nil {
			return nil, err
		}

		return users, nil
	} else if err != nil {
		return nil, err
	}

	// Если данные найдены в кеше, десериализуем их и возвращаем
	var users []*model.User
	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UserByID is the resolver for the userById field.
func (r *queryResolver) UserByID(ctx context.Context, id int) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Пытаемся получить данные из кеша Redis
	val, err := r.Models.Users.Redis.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		// Данные не найдены в кеше, загружаем из базы данных
		user, err := r.Models.Users.Get(id)
		if err != nil {
			return nil, err
		}

		// Сохраняем данные в кеш Redis
		data, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		err = r.Models.Users.Redis.Set(ctx, cacheKey, data, 10*time.Minute).Err()
		if err != nil {
			return nil, err
		}

		return user, nil
	} else if err != nil {
		return nil, err
	}

	// Если данные найдены в кеше, десериализуем их и возвращаем
	var user *model.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
