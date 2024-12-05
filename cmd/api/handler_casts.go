package main

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) createCastHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		MovieID  int64   `json:"movie_id"`
		PersonID int64   `json:"person_id"`
		JobID    int64   `json:"job_id"`
		Role     *string `json:"role"`
		Position int32   `json:"position"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cast := database.Cast{
		MovieID:  input.MovieID,
		PersonID: input.PersonID,
		JobID:    input.JobID,
		Position: input.Position,
	}

	v := validator.New()
	if input.Role != nil {
		cast.Role = *input.Role
	} else {
		v.AddError("role", "role must be provided")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	database.ValidateCast(v, &cast)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Casts.Insert(&cast)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/movies/%d", cast.MovieID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"casts": cast}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) getCastsByMovieIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	casts, err := app.models.Casts.GetByMovieID(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"casts": casts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getCastsByPersonIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	casts, err := app.models.Casts.GetByPersonID(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"casts": casts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateCastHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		MovieID  *int64  `json:"movie_id"`
		PersonID *int64  `json:"person_id"`
		JobID    *int64  `json:"job_id"`
		Role     *string `json:"role"`
		Position *int32  `json:"position"`
		Version  int32   `json:"version"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cast, err := app.models.Casts.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.MovieID != nil {
		cast.MovieID = *input.MovieID
	}
	if input.PersonID != nil {
		cast.PersonID = *input.PersonID
	}
	if input.JobID != nil {
		cast.JobID = *input.JobID
	}
	if input.Role != nil {
		cast.Role = *input.Role
	}
	if input.Position != nil {
		cast.Position = *input.Position
	}

	v := validator.New()
	database.ValidateCast(v, cast)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Casts.Update(cast)
	if err != nil {
		if errors.Is(err, database.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"cast": cast}, nil)

}

func (app *application) deleteCastHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Casts.Delete(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "cast successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
