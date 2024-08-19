package graph

import (
	"context"
	"errors"
	"github.com/olzzhas/narxozer/graph/model"
	"golang.org/x/crypto/bcrypt"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthPayload, error) {
	// Проверка, существует ли пользователь с данным email
	existingUser, err := r.Models.Users.GetByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Создание нового пользователя
	user := &model.User{
		Email:        input.Email,
		Name:         input.Name,
		Lastname:     input.Lastname,
		PasswordHash: string(hashedPassword),
		Role:         "STUDENT",
	}

	// Сохранение пользователя в базе данных
	err = r.Models.Users.Insert(user)
	if err != nil {
		return nil, err
	}

	// Генерация JWT токенов
	accessToken, err := r.JWTManager.Generate(int64(user.ID), user.Role.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := r.JWTManager.GenerateRefresh(int64(user.ID), user.Role.String())
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	// Найти пользователя в базе данных по email
	user, err := r.Models.Users.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	// Генерация токенов
	accessToken, err := r.JWTManager.Generate(int64(user.ID), user.Role.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := r.JWTManager.GenerateRefresh(int64(user.ID), user.Role.String())
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthPayload, error) {
	// Верификация Refresh токена
	claims, err := r.JWTManager.Verify(refreshToken)
	if err != nil {
		return nil, err
	}

	// Генерация нового Access токена
	accessToken, err := r.JWTManager.Generate(claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Используем старый Refresh токен
	}, nil
}
