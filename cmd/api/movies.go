package main

import (
	"database/sql"
	"errors"
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

	movie, err := app.models.Movies.Get(movieID)

	if err != nil {

		if errors.Is(err, models.ErrResourceNotFound) {
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

	requestedMovie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrResourceNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	input := struct {
		Title   *string         `json:"title"`
		Year    *int32          `json:"year"`
		Runtime *models.Runtime `json:"runtime"`
		Genres  []string        `json:"genres"`
	}{}

	err = readJSON(w, r, &input)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// if partial request is sent, update the struct partially
	if input.Title != nil {
		requestedMovie.Title = *input.Title
	}

	if input.Year != nil {
		requestedMovie.Title = *input.Title
	}

	if input.Runtime != nil {
		requestedMovie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		requestedMovie.Genres = input.Genres
	}

	validationErrors := requestedMovie.Validate()

	if len(validationErrors) != 0 {
		app.failedValidationResponse(w, r, validationErrors)
		return
	}

	err = app.models.Movies.Update(requestedMovie)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}

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

		switch {

		case errors.Is(err, models.ErrResourceNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)

		}

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
