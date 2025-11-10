package main

import (
	"fmt"
	"net/http"

	"github.com/ariffil/greenlight/internal/models"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input models.Movie

	err := readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	validationErrors := input.Validate()

	if validationErrors != nil {
		app.failedValidationResponse(w, r, validationErrors)
		return
	}

	err = app.models.Movies.Insert(&input)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", input.Id))

	response := envelope{"status": "success"}

	err = app.writeJSON(w, http.StatusCreated, response, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	movieID, err := app.readIDParam(r)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	movie, err := app.models.Movies.Get(int(movieID))

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	movieEnvelope := envelope{"movie": movie}

	err = app.writeJSON(w, http.StatusOK, movieEnvelope, nil)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return

	}

}
