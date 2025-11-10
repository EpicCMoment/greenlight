package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config      *config
	infoLogger  *log.Logger
	errorLogger *log.Logger
	db          *sql.DB
}

func main() {

	var cfg config

	app := application{
		config: &cfg,
	}

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

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
