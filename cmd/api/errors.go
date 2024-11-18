package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// var (
// 	ErrRecordNotFound = errors.New("record not found")
// 	ErrEditConflict   = errors.New("edit conflict")
// )

func (app *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := `the server encountered a problem and could not process your request`
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := `the requested resource could not be found`
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) isCtxTimeoutError(ctx context.Context, err error) bool {
	//->this fucntion returns true if the error is caused by the client/end user cancelling the http connection

	//pg returns the error "pq: canceling statement due to user request" in two cases.
	//1. the the timeout of the context has been reached, which is an error case.
	//2. when the client/end user cancels their http connection to the server. In case 2 we don't want to do anything more just return.
	//By calling ctx.Err() we get "context.Canceled" or "context.DeadlineExceeded" back. The latter we want to handle as an error.
	//Let the serverErrorResponse handle the context.DeadlineExceeded and simply return if this function returns true.
	if err.Error() == "pq: canceling statement due to user request" {
		err = ctx.Err()
		return errors.Is(err, context.Canceled)
	}
	return false
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
