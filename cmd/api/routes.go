package main

import (
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

	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)

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