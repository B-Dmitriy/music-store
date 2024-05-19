package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/B-Dmitriy/test-api/internal/config"
	"github.com/B-Dmitriy/test-api/internal/storage/sqlite"

	srv "github.com/B-Dmitriy/test-api/internal/server/http"
	lgr "github.com/B-Dmitriy/test-api/pgk/logger"
)

type App struct {
	config  *config.Config
	storage *sql.DB
	logger  *slog.Logger
	router  *http.ServeMux
}

type TestData struct {
	ID   int
	Name string
}

func New() (*App, error) {
	cfg := config.MustReadConfig()

	logger := lgr.New(cfg.Env)
	logger.Info("logger initialized", slog.String("env", cfg.Env))

	db, err := sqlite.New()
	if err != nil {
		logger.Error("storage initialization error", slog.String("text", err.Error()))
		os.Exit(1)
	}
	logger.Info("storage initialized", slog.String("driver", "sqlite3"))

	rows, err := db.Query("SELECT * FROM test;")
	if err != nil {
		logger.Error("storage initialization error", slog.String("text", err.Error()))
	}
	testData := make([]TestData, 0)

	for rows.Next() {
		item := TestData{}
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			logger.Error("database error", slog.String("text", err.Error()))
		}
		testData = append(testData, item)
	}

	fmt.Printf("%v\n", testData)

	router := srv.New(logger, db)

	return &App{
		storage: db,
		config:  cfg,
		logger:  logger,
		router:  router,
	}, nil
}

func (a *App) Run() {
	fmt.Printf("%v\n", a.config.Port)
	fmt.Printf("%V\n", a.config.Port)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", a.config.Port), a.router); err != nil {
		a.logger.Error(fmt.Sprintf("application start error: %s", err.Error()))
		os.Exit(1)
	}
}

func (a *App) Stop() {
	err := a.storage.Close()
	if err != nil {
		a.logger.Error(fmt.Sprintf("application stopped with error: %s", err.Error()))
		os.Exit(1)
	}

	a.logger.Info("application stopped success")
}
