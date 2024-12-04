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

func (app *application) getMovieKeywordsHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movieKeywords, err := app.models.CategoryItems.Get(movieId, "movie_keywords")
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie_keywords": movieKeywords}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *application) getMovieCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movieCategories, err := app.models.CategoryItems.Get(movieId, "movie_categories")
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie_categories": movieCategories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteMovieKeywordHandler(w http.ResponseWriter, r *http.Request) {
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

	err = app.models.CategoryItems.Delete(movieKeyword.MovieId, movieKeyword.CategoryId, "movie_keywords")
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie_keyword successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) deleteMovieCategoryHandler(w http.ResponseWriter, r *http.Request) {
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

	err = app.models.CategoryItems.Delete(movieKeyword.MovieId, movieKeyword.CategoryId, "movie_categories")
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie_category successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
