package storage

import "github.com/B-Dmitriy/music-store/internal/models"

type Products interface {
	GetAll() ([]models.Product, error)
}

type Categories interface {
	GetAll() ([]models.Product, error)
}

type Users interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.RegistrationData) error
}

type Tokens interface {
	GetByUserID(userID int) (*models.RefreshToken, error)
	Create(userID int, token string) error
	ChangeToken(userID int, token string) error
	RemoveByUserID(userID int) error
}
