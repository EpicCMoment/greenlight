package main

import (
	"net/http"
	"time"

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

	app.infoLogger.Printf("%+v\n", input)

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	movieID, err := app.readIDParam(r)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	movie := models.Movie{
		Id:        movieID,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	movieEnvelope := envelope{"movie": movie}

	err = app.writeJSON(w, http.StatusOK, movieEnvelope, nil)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return

	}

}
