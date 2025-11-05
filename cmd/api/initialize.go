package main

import (
	"log"
	"os"
)

func (app *application) initializeLoggers() {

	errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Llongfile)
	infoLogger := log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)

	app.errorLogger = errorLogger
	app.infoLogger = infoLogger

}
