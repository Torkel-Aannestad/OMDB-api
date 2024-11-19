package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/auth"
	"github.com/Torkel-Aannestad/MovieMaze/internal/data"
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
	v := validator.New()
	data.ValidatePasswordPlaintext(v, input.Password)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	pw_hash, err := auth.GenerateHashFromPlaintext(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	createUserParams := database.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: pw_hash,
		Activated:    false,
	}

	data.ValidateCreateUserParams(v, &createUserParams)
	valid = v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	user, err := app.model.CreateUser(ctx, createUserParams)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	token, err := auth.GenerateToken(user.ID, 3*24*time.Hour, auth.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	var userReponse = struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Activated bool   `json:"activated"`
	}{
		Name:      user.Name,
		Email:     user.Email,
		Activated: user.Activated,
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": userReponse}, nil)
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
	auth.ValidateTokenPlaintext(v, input.TokenPlaintext)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	getForTokenParams := database.GetForTokenParams{
		Hash:   []byte(input.TokenPlaintext),
		Scope:  auth.ScopeActivation,
		Expiry: time.Now(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	user, err := app.model.GetForToken(ctx, getForTokenParams)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	updateUserParams := database.UpdateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Activated:    user.Activated,
		ID:           user.ID,
	}
	_, err = app.model.UpdateUser(ctx, updateUserParams)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	deleteAllForUserParams := database.DeleteAllForUserParams{
		Scope:  auth.ScopeActivation,
		UserID: user.ID,
	}
	err = app.model.DeleteAllForUser(ctx, deleteAllForUserParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var userReponse = struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Activated bool   `json:"activated"`
	}{
		Name:      user.Name,
		Email:     user.Email,
		Activated: user.Activated,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": userReponse}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
