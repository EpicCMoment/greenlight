package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// `destination` should be a non nil pointer
func readJSON(w http.ResponseWriter, r *http.Request, destination any) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.DisallowUnknownFields()

	err := bodyDecoder.Decode(destination)

	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):

			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err

		}

	}

	hasMore := bodyDecoder.More()

	if hasMore {
		return errors.New("each request should contain a single JSON")
	}

	return nil

}
