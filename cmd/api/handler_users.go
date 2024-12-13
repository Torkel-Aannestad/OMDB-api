package main

import (
	"errors"

	"net/http"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/auth"
	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := database.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	v := validator.New()
	auth.ValidatePlaintextPassword(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	hash, err := auth.GenerateHashFromPlaintext(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.PasswordHash = hash

	database.ValidateUser(v, &user)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.models.Permissions.AddForUser(user.ID, "movies:read", "people:read", "casts:read", "jobs:read", "categories:read", "category-items:read", "movie-links:read", "people-links:read", "trailers:read", "images:write", "movies:write", "people:write", "casts:write", "jobs:write", "categories:write", "category-items:write", "movie-links:write", "people-links:write", "trailers:write", "images:write")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, database.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.backgroundJob(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var Input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &Input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	database.ValidateTokenPlaintext(v, Input.TokenPlaintext)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(database.ScopeActivation, Input.TokenPlaintext)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		} else {

			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	err = app.models.Users.Update(user)
	if err != nil {
		if errors.Is(err, database.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(database.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) addUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetById(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, "movies:read", "people:read", "casts:read", "jobs:read", "categories:read", "category-items:read", "movie-links:read", "people-links:read", "trailers:read", "images:write", "movies:write", "people:write", "casts:write", "jobs:write", "categories:write", "category-items:write", "movie-links:write", "people-links:write", "trailers:write", "images:write")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	userPermissions, err := app.models.Permissions.GetAllForUser(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"permissions": userPermissions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
