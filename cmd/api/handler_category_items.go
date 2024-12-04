package main

import (
	"errors"

	"net/http"

	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
)

func (app *application) createMovieKeywordsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		MovieId    int64 `json:"movie_id"`
		CategoryId int64 `json:"category_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movieKeyword := database.CategoryItem{
		MovieId:    input.MovieId,
		CategoryId: input.CategoryId,
	}

	err = app.models.CategoryItems.Insert(&movieKeyword, "movie_keywords")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie_keywords": movieKeyword}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *application) createMovieCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		MovieId    int64 `json:"movie_id"`
		CategoryId int64 `json:"category_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movieCategory := database.CategoryItem{
		MovieId:    input.MovieId,
		CategoryId: input.CategoryId,
	}

	err = app.models.CategoryItems.Insert(&movieCategory, "movie_categories")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie_categories": movieCategory}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
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
