package main

import (
	"fmt"
	"net/http"

	"github.com/Torkel-Aannestad/MovieMaze/internal/data"
	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "listMoviesHandler")
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

	movieParams := database.CreateMovieParams{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	data.ValidateMovie(v, &movieParams)
	valid := v.Valid()
	if !valid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movie, err := app.model.CreateMovie(r.Context(), movieParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

	movie, err := app.model.GetMovieById(r.Context(), id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) editMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "editMovieHandler")
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "deleteMovieHandler")
}
