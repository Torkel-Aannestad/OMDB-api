package main

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) createMovieLinkHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Source   string `json:"source"`
		Key      string `json:"key"`
		MovieID  int64  `json:"movie_id"`
		Language string `json:"language"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movieLink := database.MovieLink{
		Source:   input.Source,
		Key:      input.Key,
		MovieID:  input.MovieID,
		Language: input.Language,
	}

	v := validator.New()
	database.ValidateMovieLink(v, &movieLink)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.MovieLinks.Insert(&movieLink)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/categories/%d", category.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) getMovielinksHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movieLinks, err := app.models.MovieLinks.GetMovieLinks(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie_links": movieLinks}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Categories.Delete(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "category successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
