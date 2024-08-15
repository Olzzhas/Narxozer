package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/olzzhas/narxoz-sport-hub/internal/data"
	"github.com/olzzhas/narxoz-sport-hub/internal/image"
	"github.com/olzzhas/narxoz-sport-hub/internal/validator"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Surname:   input.Surname,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		_data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", _data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// TODO
	user, err := app.models.Users.GetForToken(input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
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
		return
	}

	isPassEquals, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !isPassEquals {
		app.invalidCredentialsResponse(w, r)
		return
	}

	tokens, err := app.models.AuthorizationTokens.New(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	err = app.writeJSON(w, http.StatusOK, envelope{"response": map[string]any{
		"tokens": tokens,
		"user":   user,
	}}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		fmt.Println("Error while getting refresh cookie:", err)
		return
	}

	refreshToken := cookie.Value

	fmt.Println(refreshToken)

	err = app.models.AuthorizationTokens.DeleteByToken(refreshToken)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	refreshCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshCookie)

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	var input struct {
		Name      *string `json:"name"`
		Surname   *string `json:"surname"`
		Password  *string `json:"password"`
		Activated *bool   `json:"activated"`
		ImageUrl  *string `json:"image_url"`
		Role      *string `json:"role"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Role != nil {
		user.Role = *input.Role
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Surname != nil {
		user.Surname = *input.Surname
	}

	if input.ImageUrl != nil {
		user.ImageUrl = *input.ImageUrl
	}

	if input.Password != nil {
		err = user.Password.Set(*input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if input.Activated != nil {
		user.Activated = *input.Activated
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) setProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	maxBytes := int64(8 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes) // 4 MB

	imageFile, header, err := r.FormFile("imageFile")
	if err != nil {
		switch {
		case err.Error() == "http: request body too large":
			app.badRequestResponse(w, r, fmt.Errorf("max request body size is %v bytes\n", maxBytes))
		default:
			app.logger.PrintError(fmt.Errorf("unable parse multipart/form-data: %+v", err), nil)
		}
		return
	}

	mimeType := header.Header.Get("Content-Type")

	if valid := image.IsAllowedImageType(mimeType); !valid {
		app.badRequestResponse(w, r, fmt.Errorf("image is not an allowable mime-type(png/jpeg)"))
		return
	}

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	objName, err := app.objNameFromURL(user.ImageUrl)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//imageFile, err := header.Open()
	//if err != nil {
	//	app.logger.PrintError(fmt.Errorf("failed to open image file; %v\n", err), nil)
	//	return
	//}

	imageURL, err := app.storages.ProfileImage.UpdateProfile(context.Background(), objName, imageFile)
	if err != nil {
		app.logger.PrintError(fmt.Errorf("unable to upload image to cloud provider: %v\n", err), nil)
		return
	}

	user.ImageUrl = imageURL

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
