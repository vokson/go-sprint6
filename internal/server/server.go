package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

type Server struct {
	Logger *log.Logger
	Server *http.Server
}

func InitApp(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", handlers.HandleUpload)
	mux.HandleFunc("GET /{$}", handlers.HandleMain)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{Logger: logger, Server: server}
}
