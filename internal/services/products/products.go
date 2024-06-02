package products

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/B-Dmitriy/music-store/internal/models"
	"github.com/B-Dmitriy/music-store/internal/storage/products"
	"github.com/B-Dmitriy/music-store/internal/storage/productscategories"
	"github.com/B-Dmitriy/music-store/pgk/web"
	"github.com/go-playground/validator/v10"
)

type ProductsService struct {
	logger                *slog.Logger
	productsStorage       *products.ProductsStorage
	prodCategoriesStorage *productscategories.ProductsCategoryStorage
	validator             *validator.Validate
}

func New(
	logger *slog.Logger,
	storage *products.ProductsStorage,
	pcs *productscategories.ProductsCategoryStorage,
	v *validator.Validate,
) *ProductsService {
	return &ProductsService{
		productsStorage:       storage,
		logger:                logger,
		prodCategoriesStorage: pcs,
		validator:             v,
	}
}

func (h *ProductsService) panicRecover(w http.ResponseWriter, op string) {
	if r := recover(); r != nil {
		h.logger.Error("panic in services.products", slog.String("op", op))
		web.WriteServerError(w, fmt.Errorf("server error"))
		return
	}
}

func (p *ProductsService) GetProductsList(w http.ResponseWriter, r *http.Request) {
	op := "services.products.GetProductsListHandler"
	logger := p.logger.With(slog.String("op", op))

	defer p.panicRecover(w, op)
	defer r.Body.Close()

	body := models.ProductsFilter{
		CategoryID: 0,
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerError(w, err)
		return
	}

	queryParams, err := readGetProductsQuery(r)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	switch {
	case body.CategoryID > 0:
		productsList, err := p.productsStorage.GetAllWithCategory(queryParams, body.CategoryID)
		if err != nil {
			logger.Error("storage error", slog.String("text", err.Error()))
			web.WriteServerError(w, err)
			return
		}
		web.WriteJSON(w, productsList)
		return
	default:
		productsList, err := p.productsStorage.GetAll(queryParams)
		if err != nil {
			logger.Error("storage error", slog.String("text", err.Error()))
			web.WriteServerError(w, err)
			return
		}
		web.WriteJSON(w, productsList)
		return
	}
}

func (p *ProductsService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	op := "services.categories.CreateProduct"
	logger := p.logger.With(slog.String("op", op))

	defer p.panicRecover(w, op)
	defer r.Body.Close()

	var body models.ProductCreateBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	err = p.validator.Struct(&body)
	if err != nil {
		web.WriteBadRequest(w, err.(validator.ValidationErrors))
		return
	}

	err = p.prodCategoriesStorage.CreateProductWithCategories(&body)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: products.title") {
			web.WriteBadRequest(w, fmt.Errorf("title must be unique"))
			return
		}
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
}

func (p *ProductsService) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	op := "services.categories.UpdateProduct"
	logger := p.logger.With(slog.String("op", op))

	defer p.panicRecover(w, op)
	defer r.Body.Close()

	idFromURL := r.PathValue("id")
	productID, err := strconv.Atoi(idFromURL)
	if err != nil || productID < 0 {
		web.WriteBadRequest(w, fmt.Errorf("productID must been positive int"))
		return
	}

	var body models.ProductUpdateBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	err = p.validator.Struct(&body)
	if err != nil {
		web.WriteBadRequest(w, err.(validator.ValidationErrors))
		return
	}

	err = p.prodCategoriesStorage.UpdateProductWithCategories(productID, &body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			web.WriteNotFound(w, fmt.Errorf("product not found"))
			return
		}
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
}

func (p *ProductsService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	op := "services.categories.DeleteProduct"
	logger := p.logger.With(slog.String("op", op))

	defer p.panicRecover(w, op)

	idFromURL := r.PathValue("id")
	productID, err := strconv.Atoi(idFromURL)
	if err != nil || productID < 0 {
		web.WriteBadRequest(w, fmt.Errorf("productID must been positive int"))
		return
	}

	err = p.productsStorage.Delete(productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			web.WriteNotFound(w, fmt.Errorf("product not found"))
			return
		}
		logger.Error("storage error", slog.String("text", err.Error()))
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
}

func readGetProductsQuery(r *http.Request) (*models.GetAllProductsParams, error) {
	params := models.GetAllProductsParams{
		Page:  1,
		Limit: 10,
	}

	queryParams, _ := url.ParseQuery(r.URL.RawQuery)

	if len(queryParams["page"]) != 0 {
		i, err := strconv.Atoi(queryParams["page"][0])
		if err != nil || i < 1 {
			return nil, fmt.Errorf("page must be positive integer greater than zero")
		}
		params.Page = i
	}

	if len(queryParams["limit"]) != 0 {
		i, err := strconv.Atoi(queryParams["limit"][0])
		if err != nil || i < 1 {
			return nil, fmt.Errorf("limit must be positive integer greater than zero")
		}
		params.Limit = i
	}

	return &params, nil
}
