package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ariffil/greenlight/internal/validator"
	"github.com/lib/pq"
)

type MovieModel struct {
	DB *sql.DB
}

type Movie struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int       `json:"version"`
}

func (m MovieModel) Insert(movie *Movie) error {

	stmt := `INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, version`

	row := m.DB.QueryRow(stmt, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres))

	err := row.Scan(&movie.Id, &movie.CreatedAt, &movie.Version)

	if err != nil {
		return err
	}

	return nil

}

func (m MovieModel) Delete(id int) error {
	return nil
}

func (m MovieModel) Get(id int) (*Movie, error) {

	if id < 0 {
		return nil, errors.New("id can't be less than 0")
	}

	stmt := `SELECT * FROM movies WHERE $1 = id`

	rows := m.DB.QueryRow(stmt, id)

	var data Movie

	err := rows.Scan(&data.Id, &data.CreatedAt, &data.Title, &data.Year, &data.Runtime, pq.Array(&data.Genres), &data.Version)

	if err != nil {
		return nil, err
	}

	return &data, nil

}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// if valid: returns `true`
// else `false`
func (m *Movie) Validate() (erors map[string]string) {

	v := validator.New()

	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) < 500, "title", "must not be more than 500 bytes long")

	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= 1888, "year", "must be greater than 1888")
	v.Check(m.Year <= time.Now().Year(), "year", "must not be in the future")

	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")

	if v.IsValid() {
		return nil
	} else {
		return v.Errors
	}

}
