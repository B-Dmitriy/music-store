package productscategories

import (
	"database/sql"

	"github.com/B-Dmitriy/music-store/internal/models"
)

type ProductsCategoryStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *ProductsCategoryStorage {
	return &ProductsCategoryStorage{
		db: db,
	}
}

func (pc *ProductsCategoryStorage) CreateProductWithCategories(body *models.ProductCreateBody) error {

	tx, err := pc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO products(title, price) VALUES (?, ?) RETURNING id;", body.Title, body.Price)
	if err != nil {
		return err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	queryString := "INSERT INTO lnk_category_porduct(category_id,product_id) VALUES "
	values := make([]any, 0)
	totalCategories := len(body.CategoryIDs)

	for i := range body.CategoryIDs {
		values = append(values, body.CategoryIDs[i])
		values = append(values, int(productID))

		if (i + 1) == totalCategories {
			queryString = queryString + "(?, ?);"
		} else {
			queryString = queryString + "(?, ?),"
		}
	}

	_, err = tx.Exec(queryString, values...)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (pc *ProductsCategoryStorage) UpdateProductWithCategories(productID int, body *models.ProductUpdateBody) error {

	tx, err := pc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("UPDATE products SET title=?,price=? WHERE id = ?;", body.Title, body.Price, productID)
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

	_, err = tx.Exec("DELETE from lnk_category_porduct WHERE product_id = ?;", productID)
	if err != nil {
		return err
	}

	queryString := "INSERT INTO lnk_category_porduct(category_id,product_id) VALUES "
	values := make([]any, 0)
	totalCategories := len(body.CategoryIDs)

	for i := range body.CategoryIDs {
		values = append(values, body.CategoryIDs[i])
		values = append(values, productID)

		if (i + 1) == totalCategories {
			queryString = queryString + "(?, ?);"
		} else {
			queryString = queryString + "(?, ?),"
		}
	}

	_, err = tx.Exec(queryString, values...)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
