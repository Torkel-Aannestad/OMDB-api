package main

import (
	"crypto/sha256"
	"errors"
	"strings"

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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, database.ScopeActivation, database.TokenData{})
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

func (app *application) resendActionToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if user.Activated {
		err = app.writeJSON(w, http.StatusOK, envelope{"message": "user already active, see auth/authenticate endpoint"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, database.ScopeActivation, database.TokenData{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.backgroundJob(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_resend_activation.tmpl", data)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"message": "sending new activation token to use with /users/activate endpoint"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
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

func (app *application) changeEmailHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		NewEmail string `json:"email"`
	}
	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	headerParts := strings.Split(authorizationHeader, " ")
	hash := sha256.Sum256([]byte(headerParts[1]))
	tokenHash := hash[:]

	token, err := app.models.Tokens.GetByTokenHash(database.ScopeAuthentication, tokenHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !app.models.Tokens.ValidTokenAge(time.Hour*1, token) {
		app.tokenExiredResponse(w, r)
		return
	}

	v := validator.New()
	database.ValidateEmail(v, input.NewEmail)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)
	emailVerificationToken, err := app.models.Tokens.New(user.ID, time.Hour*12, database.ScopeChangeEmail, database.TokenData{"email": input.NewEmail})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.backgroundJob(func() {
		data := map[string]any{
			"verificationToken": emailVerificationToken.Plaintext,
		}

		err = app.mailer.Send(user.Email, "change-email-verification.tmpl", data)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "a verification token will be sent to your new email"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) changeEmailVerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := app.models.Users.GetForToken(database.ScopeChangeEmail, Input.TokenPlaintext)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			v.AddError("token", "invalid or expired email verification token")
			app.failedValidationResponse(w, r, v.Errors)
		} else {

			app.serverErrorResponse(w, r, err)
		}
		return
	}
	userCurrentEmail := user.Email

	inputTokenHash := sha256.Sum256([]byte(Input.TokenPlaintext))
	emailVerificationTokenHash := inputTokenHash[:]
	token, err := app.models.Tokens.GetByTokenHash(database.ScopeChangeEmail, emailVerificationTokenHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newEmail, ok := token.Data["email"].(string)
	if !ok {
		app.logger.Error("changeEmailVerification", "message", "could not assert newEmail to string from token")
		app.serverErrorResponse(w, r, err)
		return
	}
	user.Email = newEmail
	err = app.models.Users.Update(user)
	if err != nil {
		if errors.Is(err, database.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(database.ScopeChangeEmail, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.backgroundJob(func() {
		data := map[string]any{
			"newEmail": newEmail,
		}

		err = app.mailer.Send(userCurrentEmail, "email-changed.tmpl", data)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
