package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ariffil/greenlight/internal/models"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config      *config
	infoLogger  *log.Logger
	errorLogger *log.Logger
	db          *sql.DB
	models      models.Models
}

func main() {

	var cfg config
	cfg.parseCmdFlags()

	app := application{
		config: &cfg,
	}

	app.initializeLoggers()

	err := app.openDB()

	if err != nil {
		app.errorLogger.Fatalln(err.Error())
	}

	router := app.routes()

	srv := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.infoLogger.Printf("starting %s server on http://%s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()

	app.errorLogger.Fatal(err)

}
