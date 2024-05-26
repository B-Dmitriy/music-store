package categories

import (
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/music-store/internal/storage/categories"
	"github.com/B-Dmitriy/music-store/pgk/web"
)

type CategoriesService struct {
	logger  *slog.Logger
	storage *categories.CategoriesStorage
}

func New(logger *slog.Logger, storage *categories.CategoriesStorage) *CategoriesService {
	return &CategoriesService{
		logger:  logger,
		storage: storage,
	}
}

// GetCategoriesList - curl -i -X GET "http://localhost:5050/api/categories"
func (h *CategoriesService) GetCategoriesList(w http.ResponseWriter, r *http.Request) {
	op := "server.categories.GetCategoriesList"
	logger := h.logger.With(slog.String("op", op))

	categoriesList, err := h.storage.GetAll()
	if err != nil {
		logger.Error("storage initialization error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, categoriesList)
}
