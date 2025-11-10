package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"
)

func (app *application) initializeLoggers() {

	errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Llongfile)
	infoLogger := log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)

	app.errorLogger = errorLogger
	app.infoLogger = infoLogger

}

func (app *application) openDB() error {

	dsn := os.Getenv("GREENLIGHT_DB_DSN")

	if dsn == "" {
		return errors.New("database connection string is not set in GREENLIGHT_DB_DSN environment variable")
	}

	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = conn.PingContext(ctx)

	if err != nil {
		return err
	}

	app.db = conn

	app.infoLogger.Println("database connection pool established")

	return nil

}
