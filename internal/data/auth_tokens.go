package data

import (
	"context"
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type AuthorizationToken struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthorizationTokenModel struct {
	DB *sql.DB
}

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func generateAuthenticationTokens(userID int64) (string, string, error) {
	accessClaims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSecret := os.Getenv("ACCESS_SECRET")
	accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), // Expires in 30 days
			IssuedAt:  time.Now().Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSecret := os.Getenv("REFRESH_SECRET")
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (m AuthorizationTokenModel) New(userID int64) (*AuthorizationToken, error) {
	accessToken, refreshToken, err := generateAuthenticationTokens(userID)
	if err != nil {
		return nil, err
	}

	var tokens AuthorizationToken
	tokens.UserID = userID
	tokens.RefreshToken = refreshToken
	tokens.AccessToken = accessToken

	err = m.Insert(&tokens)
	return &tokens, err
}

func (m AuthorizationTokenModel) Insert(tokens *AuthorizationToken) error {
	err := m.Delete(tokens.UserID)

	query := `
        INSERT INTO authorization_tokens (user_id, refresh_token)
        VALUES ($1, $2)
        RETURNING id
    `

	args := []interface{}{tokens.UserID, tokens.RefreshToken}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, args...)
	err = row.Scan(&tokens.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m AuthorizationTokenModel) Delete(userID int64) error {
	query := `
		DELETE FROM authorization_tokens
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID)
	return err
}

func (m AuthorizationTokenModel) DeleteByToken(refreshToken string) error {
	query := `
		DELETE FROM authorization_tokens
		WHERE refresh_token = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, refreshToken)
	return err
}

func (m AuthorizationTokenModel) Get(userID int64) (string, error) {
	query := `
		SELECT refresh_token from authorization_tokens
		WHERE authorization_tokens.user_id = $1
	`

	var refreshToken string

	err := m.DB.QueryRow(query, userID).Scan(&refreshToken)
	return refreshToken, err
}
