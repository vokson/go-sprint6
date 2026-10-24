package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := server.InitApp(logger)

	if err := app.Server.ListenAndServe(); err != nil {
		app.Logger.Fatalf("Error on server's start-up: %v\n", err)
	}
}
