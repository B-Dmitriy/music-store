package http

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/music-store/internal/services/auth"
	"github.com/B-Dmitriy/music-store/internal/services/categories"
	"github.com/B-Dmitriy/music-store/internal/services/products"
	"github.com/B-Dmitriy/music-store/pgk/password"
	"github.com/B-Dmitriy/music-store/pgk/tokens"

	categoriesStore "github.com/B-Dmitriy/music-store/internal/storage/categories"
	productsStore "github.com/B-Dmitriy/music-store/internal/storage/products"
	tokensStore "github.com/B-Dmitriy/music-store/internal/storage/tokens"
	usersStore "github.com/B-Dmitriy/music-store/internal/storage/users"
)

func New(logger *slog.Logger, db *sql.DB, pm *password.PasswordManager, tm *tokens.TokensManager) *http.ServeMux {
	server := http.NewServeMux()

	usersStorage := usersStore.New(db)
	tokensStorage := tokensStore.New(db)
	productsStorage := productsStore.New(db)
	categoriesStorage := categoriesStore.New(db)

	productsService := products.New(logger, productsStorage)
	categoriesService := categories.New(logger, categoriesStorage)
	authService := auth.New(logger, pm, tm, usersStorage, tokensStorage)

	server.HandleFunc("POST /api/login", authService.Login)
	server.HandleFunc("POST /api/logout", authService.Logout)
	server.HandleFunc("POST /api/registration", authService.Registration)
	server.HandleFunc("GET /api/products", productsService.GetProductsList)
	server.HandleFunc("GET /api/categories", categoriesService.GetCategoriesList)

	logger.Info("server routes initialization success")
	return server
}
