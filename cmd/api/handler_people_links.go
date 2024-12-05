package main

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

func (app *application) createPeopleLinkHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Source   string `json:"source"`
		Key      string `json:"key"`
		PersonID int64  `json:"person_id"`
		Language string `json:"language"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	peopleLink := database.PeopleLink{
		Source:   input.Source,
		Key:      input.Key,
		PersonID: input.PersonID,
		Language: input.Language,
	}

	v := validator.New()
	database.ValidatePeopleLink(v, &peopleLink)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.PeopleLinks.Insert(&peopleLink)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/people-links/%d", peopleLink.PersonID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"people_links": peopleLink}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) getPeopleLinksHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	peopleLinks, err := app.models.PeopleLinks.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"people_links": peopleLinks}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deletePeopleLinkHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.PeopleLinks.Delete(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "peoplelink successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
