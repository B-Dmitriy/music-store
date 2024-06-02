package products

import (
	"database/sql"

	"github.com/B-Dmitriy/music-store/internal/models"
)

type ProductsStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *ProductsStorage {
	return &ProductsStorage{
		db: db,
	}
}

func (s *ProductsStorage) GetAllWithCategory(params *models.GetAllProductsParams, categoryID int) ([]models.Product, error) {
	products := make([]models.Product, 0)
	offset := (params.Page - 1) * params.Limit

	rows, err := s.db.Query("SELECT * FROM products as p WHERE p.id IN (SELECT product_id FROM lnk_category_porduct WHERE category_id = ?) LIMIT ? OFFSET ?;", categoryID, params.Limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := models.Product{}
		if err := rows.Scan(&item.ID, &item.Title, &item.Price); err != nil {

			return nil, err
		}
		products = append(products, item)
	}

	return products, nil
}

func (s *ProductsStorage) GetAll(params *models.GetAllProductsParams) ([]models.Product, error) {
	products := make([]models.Product, 0)
	offset := (params.Page - 1) * params.Limit

	rows, err := s.db.Query("SELECT * FROM products LIMIT ? OFFSET ?;", params.Limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := models.Product{}
		if err := rows.Scan(&item.ID, &item.Title, &item.Price); err != nil {

			return nil, err
		}
		products = append(products, item)
	}

	return products, nil
}

func (s *ProductsStorage) Delete(productID int) error {
	result, err := s.db.Exec("DELETE FROM products WHERE id = ?;", productID)
	if err != nil {
		return err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return sql.ErrNoRows
	}

	return nil
}
