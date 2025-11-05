package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

// this function don't really need to be a method of *application
// but it doesn't hurt to be that way
func (app *application) readIDParam(r *http.Request) (int64, error) {

	if r == nil {
		return 0, errors.New("can't extract the movie ID from a nil request")
	}

	params := httprouter.ParamsFromContext(r.Context())

	movieID, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || movieID < 0 {

		return 0, errors.New("invalid movie ID")
	}

	return movieID, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, dt envelope, headers http.Header) error {

	js, err := json.MarshalIndent(dt, "", "\t")

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil

}
