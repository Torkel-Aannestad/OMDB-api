package main

import (
	"errors"
	"fmt"
	"time"

	"net/http"

	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) createPeopleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string     `json:"name"`
		Birthday time.Time  `json:"birthday,omitempty"`
		Deathday *time.Time `json:"deathday,omitempty"`
		Gender   string     `json:"gender,omitempty"`
		Aliases  *[]string  `json:"aliases,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	person := database.Person{
		Name:     input.Name,
		Birthday: input.Birthday,
		Gender:   input.Gender,
	}

	if input.Deathday != nil {
		person.Deathday = *input.Deathday
	}
	if input.Aliases != nil {
		person.Aliases = *input.Aliases
	}

	v := validator.New()
	database.ValidatePeople(v, &person)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.People.Insert(&person)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/people/%d", person.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"people": person}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) getpeopleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	person, err := app.models.People.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"person": person}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listPeopleHandler(w http.ResponseWriter, r *http.Request) {
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
	input.Filters.SortSafelist = []string{"id", "name", "birthday", "-id", "-name", "-birthday"}

	database.ValidateFilters(v, input.Filters)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	people, metadata, err := app.models.People.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"people": people, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePeopleHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) deletePeopleHandler(w http.ResponseWriter, r *http.Request) {
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
