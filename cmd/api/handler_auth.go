package main

import (
	"fmt"
	"net/http"
)

func (app *application) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Hi here")
}
