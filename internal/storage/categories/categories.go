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

func (s *CategoriesStorage) Create(name string) error {
	_, err := s.db.Exec("INSERT INTO categories (name) VALUES (?);", name)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoriesStorage) Update(categoryID int, name string) error {
	result, err := s.db.Exec("UPDATE categories SET name=? WHERE id = ?;", name, categoryID)
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

func (s *CategoriesStorage) Delete(categoryID int) error {
	result, err := s.db.Exec("DELETE FROM categories WHERE id = ?;", categoryID)
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
