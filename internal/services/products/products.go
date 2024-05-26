package products

import (
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/music-store/internal/storage/products"
	"github.com/B-Dmitriy/music-store/pgk/web"
)

type ProductsService struct {
	logger  *slog.Logger
	storage *products.ProductsStorage
}

func New(logger *slog.Logger, storage *products.ProductsStorage) *ProductsService {
	return &ProductsService{
		storage: storage,
		logger:  logger,
	}
}

// GetProductsList - curl -i -X GET "http://localhost:5050/api/products"
func (h *ProductsService) GetProductsList(w http.ResponseWriter, r *http.Request) {
	op := "server.products.GetProductsListHandler"
	logger := h.logger.With(slog.String("op", op))

	productsList, err := h.storage.GetAll()
	if err != nil {
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, productsList)
}
