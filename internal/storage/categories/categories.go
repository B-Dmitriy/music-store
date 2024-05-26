package categories

import (
	"database/sql"

	"github.com/B-Dmitriy/music-store/internal/models"
)

type CategoriesStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *CategoriesStorage {
	return &CategoriesStorage{
		db: db,
	}
}

func (s *CategoriesStorage) GetAll() ([]models.Category, error) {
	categoriesList := make([]models.Category, 0)

	rows, err := s.db.Query("SELECT * FROM categories;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := models.Category{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categoriesList = append(categoriesList, c)
	}

	return categoriesList, nil
}
