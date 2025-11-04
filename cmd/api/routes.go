package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	return router
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string {
		"status": "available",
		"environment": app.config.env,
		"version": version,
	}

	js, err := json.Marshal(data)

	// an error occured while marshalling
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")


	js = append(js, '\n')

	w.Write(js)

}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "creating a movie ...")


}
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	movieID, err := app.readIDParam(r)

	if err != nil {
		app.logger.Println(err)
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "showing the details of movie %d\n", movieID)


}