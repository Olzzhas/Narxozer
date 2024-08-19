package main

import (
	"encoding/json"
	"github.com/olzzhas/narxozer/graph/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email                 string  `json:"email"`
		Name                  string  `json:"name"`
		Lastname              string  `json:"lastname"`
		Password              string  `json:"password"`
		ImageURL              *string `json:"imageURL,omitempty"`
		AdditionalInformation *string `json:"additionalInformation,omitempty"`
		Course                *int    `json:"course,omitempty"`
		Major                 *string `json:"major,omitempty"`
		Degree                *string `json:"degree,omitempty"`
		Faculty               *string `json:"faculty,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Проверка обязательных полей
	if input.Email == "" || input.Name == "" || input.Lastname == "" || input.Password == "" {
		app.failedValidationResponse(w, r, map[string]string{
			"email":    "This field is required",
			"name":     "This field is required",
			"lastname": "This field is required",
			"password": "This field is required",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := &model.User{
		Email:                 input.Email,
		Name:                  input.Name,
		Lastname:              input.Lastname,
		PasswordHash:          string(hashedPassword),
		Role:                  model.RoleStudent,
		ImageURL:              input.ImageURL,
		AdditionalInformation: input.AdditionalInformation,
		Course:                input.Course,
		Major:                 input.Major,
		Degree:                input.Degree,
		Faculty:               input.Faculty,
		CreatedAt:             time.Now().Format(time.RFC3339),
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	accessToken, err := app.jwtManager.Generate(int64(user.ID), user.Role.String())
	refreshToken, err := app.jwtManager.GenerateRefresh(int64(user.ID), user.Role.String())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":                    user.ID,
			"email":                 user.Email,
			"name":                  user.Name,
			"lastname":              user.Lastname,
			"role":                  user.Role,
			"imageURL":              user.ImageURL,
			"additionalInformation": user.AdditionalInformation,
			"course":                user.Course,
			"createdAt":             user.CreatedAt,
			"updatedAt":             user.UpdatedAt,
			"major":                 user.Major,
			"degree":                user.Degree,
			"faculty":               user.Faculty,
		},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Email == "" || input.Password == "" {
		app.failedValidationResponse(w, r, map[string]string{
			"email":    "This field is required",
			"password": "This field is required",
		})
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if user == nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	accessToken, err := app.jwtManager.Generate(int64(user.ID), user.Role.String())
	refreshToken, err := app.jwtManager.GenerateRefresh(int64(user.ID), user.Role.String())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":                    user.ID,
			"email":                 user.Email,
			"name":                  user.Name,
			"lastname":              user.Lastname,
			"role":                  user.Role,
			"imageURL":              user.ImageURL,
			"additionalInformation": user.AdditionalInformation,
			"course":                user.Course,
			"createdAt":             user.CreatedAt,
			"updatedAt":             user.UpdatedAt,
			"major":                 user.Major,
			"degree":                user.Degree,
			"faculty":               user.Faculty,
		},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	claims, err := app.jwtManager.Verify(input.RefreshToken)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	accessToken, err := app.jwtManager.Generate(claims.UserID, claims.Role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"access_token": accessToken,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
