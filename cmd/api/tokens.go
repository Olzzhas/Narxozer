package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/olzzhas/narxozer/internal/data"
	"github.com/olzzhas/narxozer/internal/validator"
	"net/http"
	"os"
	"time"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		fmt.Println("Error while getting refresh cookie:", err)
		return
	}
	refreshToken := cookie.Value

	refreshSecret := os.Getenv("REFRESH_SECRET")

	claims, err := app.validateToken(refreshToken, refreshSecret)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	id := claims["user_id"].(float64)

	v := validator.New()

	if data.ValidateToken(v, refreshToken); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	refreshToken, err = app.models.AuthorizationTokens.Get(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			v.AddError("token", "Invalid token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if claims == nil || refreshToken == "" {
		app.invalidCredentialsResponse(w, r)
		return
	}

	tokens, err := app.models.AuthorizationTokens.New(int64(id))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &newCookie)

	err = app.writeJSON(w, http.StatusOK, envelope{"tokens": tokens}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) validateToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect method for token signing")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token format")
	}

	return claims, nil
}
