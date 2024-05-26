package tokens

import (
	"database/sql"
	
	"github.com/B-Dmitriy/music-store/internal/models"
)

type TokenStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *TokenStorage {
	return &TokenStorage{
		db: db,
	}
}

func (t *TokenStorage) Create(userID int, token string) error {
	_, err := t.db.Exec("INSERT INTO refresh_tokens (user_id, refresh_token) VALUES (?, ?)", userID, token)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenStorage) GetByUserID(userID int) (*models.RefreshToken, error) {
	var token models.RefreshToken
	row := t.db.QueryRow("SELECT * FROM refresh_tokens WHERE user_id = ?", userID)

	err := row.Scan(&token.ID, &token.UserID, &token.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *TokenStorage) ChangeToken(userID int, token string) error {
	_, err := t.db.Exec("UPDATE refresh_tokens SET refresh_token = ? WHERE user_id = ?;", token, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenStorage) RemoveByUserID(userID int) error {
	_, err := t.db.Exec("DELETE FROM refresh_tokens WHERE user_id = ?;", userID)
	if err != nil {
		return err
	}

	return nil
}
