package graph

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/olzzhas/narxozer/graph/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthPayload, error) {
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
	accessToken, err := r.JWTManager.Generate(int64(user.ID), "STUDENT")
	if err != nil {
		return nil, err
	}

	refreshToken, err := r.JWTManager.GenerateRefresh(int64(user.ID), "STUDENT")
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (manager *JWTManager) Generate(userID int64, role string) (string, error) {
	claims := &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
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

func (manager *JWTManager) GenerateRefresh(userID int64, role string) (string, error) {
	claims := &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}
