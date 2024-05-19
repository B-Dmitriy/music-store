package http

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func New(logger *slog.Logger, storage *sql.DB) *http.ServeMux {
	server := http.NewServeMux()

	server.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("/login request", slog.String("from", r.RemoteAddr))
		_, _ = w.Write([]byte("login handler"))
	})

	logger.Info("server routes initialization success")
	return server
}
