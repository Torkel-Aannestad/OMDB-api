package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Torkel-Aannestad/OMDB-api/internal/auth"
	"github.com/Torkel-Aannestad/OMDB-api/internal/database"
	"github.com/Torkel-Aannestad/OMDB-api/internal/validator"
)

func (app *application) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	database.ValidateEmail(v, input.Email)
	auth.ValidatePlaintextPassword(v, input.Password)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.invalidCredentialsResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := auth.ComparePasswordAndHash(input.Password, user.PasswordHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	authToken, err := app.models.Tokens.New(user.ID, time.Hour*24, database.ScopeAuthentication, database.TokenData{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"authentication_token": authToken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	auth.ValidatePlaintextPassword(v, input.NewPassword)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user := app.contextGetUser(r)
	match, err := auth.ComparePasswordAndHash(input.CurrentPassword, user.PasswordHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	newPasswordHash, err := auth.GenerateFromPassword(input.NewPassword, auth.DefaultParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.PasswordHash = newPasswordHash
	err = app.models.Users.Update(user)
	if err != nil {
		if errors.Is(err, database.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(database.ScopeAuthentication, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	authToken, err := app.models.Tokens.New(user.ID, time.Hour*24, database.ScopeAuthentication, database.TokenData{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.backgroundJob(func() {
		err = app.mailer.Send(user.Email, "password-changed.tmpl", nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"authentication_token": authToken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, time.Hour*1, database.ScopePasswordReset, database.TokenData{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.backgroundJob(func() {
		data := map[string]any{
			"passwordResetToken": token.Plaintext,
		}
		err = app.mailer.Send(user.Email, "password-reset.tmpl", data)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "verfification token will be sendt to your email"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) resetPasswordVerifyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
		NewPassword    string `json:"new_password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	auth.ValidatePlaintextPassword(v, input.NewPassword)
	database.ValidateTokenPlaintext(v, input.TokenPlaintext)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(database.ScopePasswordReset, input.TokenPlaintext)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	newPasswordHash, err := auth.GenerateFromPassword(input.NewPassword, auth.DefaultParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.PasswordHash = newPasswordHash
	err = app.models.Users.Update(user)
	if err != nil {
		if errors.Is(err, database.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(database.ScopeAuthentication, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	authToken, err := app.models.Tokens.New(user.ID, time.Hour*24, database.ScopeAuthentication, database.TokenData{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.backgroundJob(func() {
		err = app.mailer.Send(user.Email, "password-changed.tmpl", nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"authentication_token": authToken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := app.models.Tokens.DeleteAllForUser(database.ScopeAuthentication, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Tokens.DeleteAllForUser(database.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Tokens.DeleteAllForUser(database.ScopePasswordReset, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "session tokens revoked"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
