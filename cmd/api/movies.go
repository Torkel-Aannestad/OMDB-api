package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/data"
	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Page = app.readInt(qs, "page", 0, v)
	input.PageSize = app.readInt(qs, "pagesize", 20, v)
	input.Sort = app.readString(qs, "sort", "id")

	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int64    `json:"year"`
		Runtime int64    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	createMovieParams := database.CreateMovieParams{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	data.ValidateCreateMovieParams(v, &createMovieParams)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	movie, err := app.model.CreateMovie(ctx, createMovieParams)
	if err != nil {
		switch {
		case app.isCtxTimeoutError(ctx, err):
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	movie, err := app.model.GetMovieById(ctx, id)
	if err != nil {
		switch {
		case app.isCtxTimeoutError(ctx, err):
			return
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	movie, err := app.model.GetMovieById(ctx, id)
	if err != nil {
		switch {
		case app.isCtxTimeoutError(ctx, err):
			return
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   *string  `json:"title"`
		Year    *int64   `json:"year"`
		Runtime *int64   `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updateMovieParams := database.UpdateMovieParams{
		ID:      movie.ID,
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: movie.Runtime,
		Genres:  movie.Genres,
	}
	if input.Title != nil {
		updateMovieParams.Title = *input.Title
	}
	if input.Year != nil {
		updateMovieParams.Year = *input.Year
	}
	if input.Runtime != nil {
		updateMovieParams.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		updateMovieParams.Genres = input.Genres
	}

	v := validator.New()
	data.ValidateUpdateMovieParams(v, &updateMovieParams)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	updatedMovie, err := app.model.UpdateMovie(ctx, updateMovieParams)
	if err != nil {
		switch {
		case app.isCtxTimeoutError(ctx, err):
			return
		case errors.Is(err, sql.ErrNoRows):
			//data race condition met
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": updatedMovie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	affercedRows, err := app.model.DeleteMovie(ctx, id)
	if err != nil {
		switch {
		case app.isCtxTimeoutError(ctx, err):
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if affercedRows != 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": "movie was successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
