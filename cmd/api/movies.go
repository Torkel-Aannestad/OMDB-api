package main

import (
	"fmt"
	"net/http"
)

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "listMoviesHandler")
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "createMovieHandler")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "showMovieHandler")
}

func (app *application) editMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "editMovieHandler")
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "deleteMovieHandler")
}
