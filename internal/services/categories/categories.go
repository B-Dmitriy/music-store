package categories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/B-Dmitriy/music-store/internal/models"
	"github.com/B-Dmitriy/music-store/internal/storage/categories"
	"github.com/B-Dmitriy/music-store/pgk/web"
	"github.com/go-playground/validator/v10"
)

type CategoriesService struct {
	logger    *slog.Logger
	storage   *categories.CategoriesStorage
	validator *validator.Validate
}

func New(logger *slog.Logger, storage *categories.CategoriesStorage, v *validator.Validate) *CategoriesService {
	return &CategoriesService{
		logger:    logger,
		storage:   storage,
		validator: v,
	}
}

func (h *CategoriesService) panicRecover(w http.ResponseWriter, op string) {
	if r := recover(); r != nil {
		h.logger.Error("panic in services.categories", slog.String("op", op))
		web.WriteServerError(w, fmt.Errorf("server error"))
		return
	}
}

func (h *CategoriesService) GetCategoriesList(w http.ResponseWriter, r *http.Request) {
	op := "services.categories.GetCategoriesList"
	logger := h.logger.With(slog.String("op", op))

	defer h.panicRecover(w, op)

	categoriesList, err := h.storage.GetAll()
	if err != nil {
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, categoriesList)
}

func (h *CategoriesService) CreateCategory(w http.ResponseWriter, r *http.Request) {
	op := "services.categories.CreateCategory"
	logger := h.logger.With(slog.String("op", op))

	defer h.panicRecover(w, op)
	defer r.Body.Close()

	var body models.CreateCategoryBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	err = h.validator.Struct(&body)
	if err != nil {
		web.WriteBadRequest(w, err.(validator.ValidationErrors))
		return
	}

	err = h.storage.Create(body.Name)
	if err != nil {
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
}

func (h *CategoriesService) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	op := "services.categories.UpdateCategory"
	logger := h.logger.With(slog.String("op", op))

	defer h.panicRecover(w, op)
	defer r.Body.Close()

	idFromURL := r.PathValue("id")
	categoryID, err := strconv.Atoi(idFromURL)
	if err != nil || categoryID < 0 {
		web.WriteBadRequest(w, fmt.Errorf("categoryID must been positive int"))
		return
	}

	var body models.UpdateCategoryBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	err = h.validator.Struct(&body)
	if err != nil {
		web.WriteBadRequest(w, err.(validator.ValidationErrors))
		return
	}

	err = h.storage.Update(categoryID, body.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			web.WriteNotFound(w, fmt.Errorf("category not found"))
			return
		}
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
}

func (h *CategoriesService) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	op := "services.categories.DeleteCategory"
	logger := h.logger.With(slog.String("op", op))

	defer h.panicRecover(w, op)

	idFromURL := r.PathValue("id")
	categoryID, err := strconv.Atoi(idFromURL)
	if err != nil || categoryID < 0 {
		web.WriteBadRequest(w, fmt.Errorf("categoryID must been positive int"))
		return
	}

	err = h.storage.Delete(categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			web.WriteNotFound(w, fmt.Errorf("category not found"))
			return
		}
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
}
