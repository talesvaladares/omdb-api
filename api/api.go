package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"omdb-api/omdb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(apiKey string) http.Handler {
	routes := chi.NewMux()

	routes.Use(middleware.Recoverer)
	routes.Use(middleware.RequestID)
	routes.Use(middleware.Logger)

	routes.Get("/", handleSearchMovie(apiKey))

	return routes
}

func handleSearchMovie(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("s")
		res, err := omdb.Search(apiKey, search)

		if err != nil {
			sendJSON(w, Response{Error: "something wrong with omdb"}, http.StatusBadGateway)
			return
		}
		sendJSON(w, Response{Data: res}, http.StatusOK)
	}
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, res Response, status int) {
	data, err := json.Marshal(res)
	if err != nil {
		sendJSON(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("Error ao fazer marshal de json", "error", err)
		return
	}
}
