package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/B-Dmitriy/music-store/internal/config"
	"github.com/B-Dmitriy/music-store/internal/storage"
	"github.com/B-Dmitriy/music-store/pgk/password"
	"github.com/B-Dmitriy/music-store/pgk/tokens"

	srv "github.com/B-Dmitriy/music-store/internal/server/http"
	lgr "github.com/B-Dmitriy/music-store/pgk/logger"
)

type App struct {
	config  *config.Config
	storage *sql.DB
	logger  *slog.Logger
	router  *http.ServeMux
}

func New() (*App, error) {
	cfg := config.MustReadConfig()

	logger := lgr.New(cfg.Env)
	logger.Info("logger initialized", slog.String("env", cfg.Env))

	pm := password.New(cfg.PassCost)
	logger.Info("password manager initialized")

	tm := tokens.New(cfg.SecretKey)
	logger.Info("tokens manager initialized")

	db, err := storage.New()
	if err != nil {
		logger.Error("storage initialization error", slog.String("text", err.Error()))
		return nil, err
	}
	logger.Info("storage initialized", slog.String("driver", "sqlite3"))

	router := srv.New(logger, db, pm, tm)
	logger.Info("router initialized")

	return &App{
		storage: db,
		config:  cfg,
		logger:  logger,
		router:  router,
	}, nil
}

func (a *App) Run() {
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", a.config.Port), a.router); err != nil {
		a.logger.Error(fmt.Sprintf("application start error: %s", err.Error()))
		os.Exit(1)
	}
}
