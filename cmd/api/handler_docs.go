package main

import (
	"net/http"

	"github.com/Torkel-Aannestad/OMDB-api/assets"
)

func (app *application) getDocs(w http.ResponseWriter, r *http.Request) {
	file, err := assets.EmbededFiles.ReadFile("docs/docs.html")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(file)

}
