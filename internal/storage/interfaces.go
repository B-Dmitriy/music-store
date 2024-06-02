package storage

import "github.com/B-Dmitriy/music-store/internal/models"

type Products interface {
	GetAll(params *models.GetAllProductsParams) ([]models.Product, error)
	GetAllWithCategory(params *models.GetAllProductsParams, categoryID int) ([]models.Product, error)
	Delete(productID int) error
}

type Categories interface {
	GetAll() ([]models.Category, error)
	Create(name string) error
	Update(categoryID int, name string) error
	Delete(categoryID int) error
}

type Users interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.RegistrationData) error
}

type Tokens interface {
	GetByUserID(userID int) (*models.RefreshToken, error)
	CheckToken(userID int) (bool, error)
	Create(userID int, token string) error
	ChangeToken(userID int, token string) error
	RemoveByUserID(userID int) error
}

type ProductsCategories interface {
	CreateProductWithCategories(body *models.ProductCreateBody) error
	UpdateProductWithCategories(productID int, body *models.ProductUpdateBody) error
}
