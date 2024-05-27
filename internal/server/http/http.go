package http

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
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
	validate := validator.New()

	usersStorage := usersStore.New(db)
	tokensStorage := tokensStore.New(db)
	productsStorage := productsStore.New(db)
	categoriesStorage := categoriesStore.New(db)

	productsService := products.New(logger, productsStorage)
	categoriesService := categories.New(logger, categoriesStorage, validate)
	authService := auth.New(logger, pm, tm, usersStorage, tokensStorage)

	server.HandleFunc("POST /api/login", authService.Login)
	server.HandleFunc("POST /api/logout", authService.Logout)
	server.HandleFunc("POST /api/refresh", authService.Refresh)
	server.HandleFunc("POST /api/registration", authService.Registration)
	server.HandleFunc("GET /api/products", productsService.GetProductsList)
	server.HandleFunc("GET /api/categories", categoriesService.GetCategoriesList)

	server.Handle("POST /api/categories", authService.AuthMiddleware(http.HandlerFunc(categoriesService.CreateCategory)))
	server.Handle("PUT /api/categories/{id}", authService.AuthMiddleware(http.HandlerFunc(categoriesService.UpdateCategory)))
	server.Handle("DELETE /api/categories/{id}", authService.AuthMiddleware(http.HandlerFunc(categoriesService.DeleteCategory)))

	logger.Info("server routes initialization success")
	return server
}
