package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
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

	maxIdleTimeout, err := time.ParseDuration(app.config.db.maxIdleTime)

	if err != nil {
		return err
	}

	// set database connection pool settings
	conn.SetMaxOpenConns(app.config.db.maxOpenConns)
	conn.SetMaxIdleConns(app.config.db.maxIdleConns)
	conn.SetConnMaxIdleTime(maxIdleTimeout)

	app.db = conn
	app.models.Movies.DB = conn

	app.infoLogger.Println("database connection pool established")

	return nil

}

func (c *config) parseCmdFlags() {

	flag.IntVar(&c.port, "port", 4000, "API server port")
	flag.StringVar(&c.env, "env", "development", "Environment (development|staging|production)")

	flag.IntVar(&c.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.IntVar(&c.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.StringVar(&c.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

}
