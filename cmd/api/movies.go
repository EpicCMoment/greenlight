package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ariffil/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "creating a movie ...")

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	movieID, err := app.readIDParam(r)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	movie := data.Movie{
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
		app.serverErrorResponse(w, r, err)
		return

	}

}
