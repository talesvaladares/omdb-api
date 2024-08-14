package main

import (
	"log/slog"
	"net/http"
	"omdb-api/api"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		os.Exit(1)
	}

	slog.Info("all systems offline")
}

func run() error {
	err := godotenv.Load()

	if err != nil {
		slog.Error("Erro ao carregar o arquivo .env: %v", "error", err)
	}

	apiKey := os.Getenv("OMDB_KEY")

	handler := api.NewHandler(apiKey)

	server := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
