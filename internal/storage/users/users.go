package users

import (
	"database/sql"

	"github.com/B-Dmitriy/music-store/internal/models"
)

type UsersStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}

func (us *UsersStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := us.db.QueryRow("SELECT * FROM users WHERE email = ?", email)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RoleID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UsersStorage) CreateUser(user *models.RegistrationData) error {
	_, err := us.db.Exec("INSERT INTO users (email, password, username) VALUES (?, ?, ?)", user.Email, user.Password, user.Username)
	if err != nil {
		return err
	}

	return nil
}
