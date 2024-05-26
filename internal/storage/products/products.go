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

func (s *ProductsStorage) GetAll() ([]models.Product, error) {
	products := make([]models.Product, 0)

	rows, err := s.db.Query("SELECT * FROM products;")
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
