package main

import (
	"errors"
	"fmt"
	"time"

	"net/http"

	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string             `json:"name"`
		ParentID    database.NullInt64 `json:"parent_id,omitempty"`
		Date        time.Time          `json:"date"`
		SeriesID    database.NullInt64 `json:"series_id,omitempty"`
		Kind        string             `json:"kind"`
		Runtime     int64              `json:"runtime"`
		Budget      *float64           `json:"budget,omitempty"`
		Revenue     *float64           `json:"revenue,omitempty"`
		Homepage    *string            `json:"homepage,omitempty"`
		VoteAverage float64            `json:"vote_average"`
		VotesCount  int64              `json:"votes_count"`
		Abstract    *string            `json:"abstract,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := database.Movie{
		Name:        input.Name,
		ParentID:    input.ParentID,
		SeriesID:    input.SeriesID,
		Date:        input.Date,
		Kind:        input.Kind,
		Runtime:     input.Runtime,
		Budget:      0,
		Revenue:     0,
		Homepage:    "",
		VoteAvarage: input.VoteAverage,
		VoteCount:   input.VotesCount,
		Abstract:    "",
	}

	if input.Budget != nil {
		movie.Budget = *input.Budget
	}
	if input.Revenue != nil {
		movie.Revenue = *input.Revenue
	}
	if input.Homepage != nil {
		movie.Homepage = *input.Homepage
	}
	if input.Abstract != nil {
		movie.Abstract = *input.Abstract
	}

	v := validator.New()
	database.ValidateMovie(v, &movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(&movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		database.Filters
	}

	qs := r.URL.Query()

	v := validator.New()

	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "date", "runtime", "-id", "-title", "-date", "-runtime"}

	database.ValidateFilters(v, input.Filters)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
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

	var input struct {
		Name        *string             `json:"name"`
		ParentID    *database.NullInt64 `json:"parent_id,omitempty"`
		Date        *time.Time          `json:"date"`
		SeriesID    *database.NullInt64 `json:"series_id,omitempty"`
		Kind        *string             `json:"kind"`
		Runtime     *int64              `json:"runtime"`
		Budget      *float64            `json:"budget,omitempty"`
		Revenue     *float64            `json:"revenue,omitempty"`
		Homepage    *string             `json:"homepage,omitempty"`
		VoteAverage *float64            `json:"vote_average"`
		VotesCount  *int64              `json:"votes_count"`
		Abstract    *string             `json:"abstract,omitempty"`
		Version     *int32              `json:"version"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.Name != nil {
		movie.Name = *input.Name
	}
	if input.ParentID != nil {
		movie.ParentID = *input.ParentID
	}
	if input.Date != nil {
		movie.Date = *input.Date
	}
	if input.SeriesID != nil {
		movie.SeriesID = *input.SeriesID
	}
	if input.Kind != nil {
		movie.Kind = *input.Kind
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Budget != nil {
		movie.Budget = *input.Budget
	}
	if input.Revenue != nil {
		movie.Revenue = *input.Revenue
	}
	if input.Homepage != nil {
		movie.Homepage = *input.Homepage
	}
	if input.Abstract != nil {
		movie.Abstract = *input.Abstract
	}

	v := validator.New()
	database.ValidateMovie(v, movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		if errors.Is(err, database.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
