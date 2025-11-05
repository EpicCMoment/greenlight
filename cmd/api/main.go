package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config      config
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	app := application{
		config: cfg,
	}

	app.initializeLoggers()

	router := app.routes()

	srv := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.infoLogger.Printf("starting %s server on http://%s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()

	app.errorLogger.Fatal(err)

}
