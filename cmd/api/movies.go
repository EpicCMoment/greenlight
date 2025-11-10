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

	response := envelope{"status": "created successfully"}

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

		if err.Error() == models.ErrRecordNotFound {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	movieEnvelope := envelope{"movie": movie}

	err = app.writeJSON(w, http.StatusOK, movieEnvelope, nil)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return

	}

}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input models.Movie
	input.Id = id

	err = readJSON(w, r, &input)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Movies.Update(&input)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"status": "updated successfully",
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Movies.Delete(int(id))

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"status": "deleted successfully",
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
