package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/B-Dmitriy/music-store/internal/models"
	"github.com/B-Dmitriy/music-store/internal/storage/users"
	"github.com/B-Dmitriy/music-store/pgk/password"
	"github.com/B-Dmitriy/music-store/pgk/tokens"
	"github.com/B-Dmitriy/music-store/pgk/web"

	tokensStore "github.com/B-Dmitriy/music-store/internal/storage/tokens"
)

type AuthService struct {
	logger        *slog.Logger
	passManager   *password.PasswordManager
	tokensManager *tokens.TokensManager
	usersStorage  *users.UsersStorage
	tokensStorage *tokensStore.TokenStorage
}

func New(logger *slog.Logger, pm *password.PasswordManager, tm *tokens.TokensManager, us *users.UsersStorage, ts *tokensStore.TokenStorage) *AuthService {
	return &AuthService{
		logger:        logger,
		passManager:   pm,
		tokensManager: tm,
		usersStorage:  us,
		tokensStorage: ts,
	}
}

// Login - curl -i -X POST -d '{"email": "test2@mail.ru", "password": "qwerty123"}' http://localhost:5050/api/login
func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var loginData models.LoginData
	_ = json.NewDecoder(r.Body).Decode(&loginData)

	candidate, err := a.usersStorage.GetUserByEmail(loginData.Email)
	if err != nil {
		web.WriteBadRequest(w, fmt.Errorf("email or password error"))
		return
	}

	if !a.passManager.CheckPasswordHash(loginData.Password, candidate.Password) {
		web.WriteBadRequest(w, fmt.Errorf("email or password error"))
		return
	}

	tkns, err := a.tokensManager.GenerateJWTTokens(candidate.ID, candidate.RoleID)
	if err != nil {
		web.WriteServerError(w, err)
		return
	}

	// TODO: test and refactoring
	_, err = a.tokensStorage.GetByUserID(candidate.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = a.tokensStorage.Create(candidate.ID, tkns.RefreshToken)
			if err != nil {
				web.WriteServerError(w, err)
				return
			}
		} else {
			web.WriteServerError(w, err)
			return
		}
	} else {
		err = a.tokensStorage.ChangeToken(candidate.ID, tkns.RefreshToken)
		if err != nil {
			web.WriteServerError(w, err)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tkns.RefreshToken,
		HttpOnly: true,
		MaxAge:   14 * 24 * 3600, // 14 days
	})

	web.WriteJSON(w, &models.LoginRequest{
		Username:     candidate.Username,
		UserRoleID:   candidate.RoleID,
		AccessToken:  tkns.AcceptToken,
		RefreshToken: tkns.RefreshToken,
	})
}

func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: refactoring
	headerToken := r.Header.Get("Authorization")

	authToken := strings.Split(headerToken, " ")

	// authToken[1] if len == 0 panic
	userData, err := a.tokensManager.VerifyJWTToken(authToken[1])
	if err != nil {
		web.WriteBadRequest(w, fmt.Errorf("invalid token"))
		return
	}

	err = a.tokensStorage.RemoveByUserID(userData.UserID)
	if err != nil {
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
}

// Registration - curl -i -X POST -d '{"email": "test2@mail.ru", "password": "qwerty123", "username":"user2"}' http://localhost:5050/api/registration
func (a *AuthService) Registration(w http.ResponseWriter, r *http.Request) {
	var userData models.RegistrationData
	err := json.NewDecoder(r.Body).Decode(&userData) //TODO: закрывает ли чтение body
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	_, err = a.usersStorage.GetUserByEmail(userData.Email)
	if errors.Is(err, sql.ErrNoRows) {
		passHash, err := a.passManager.HashPassword(userData.Password)
		if err != nil {
			web.WriteServerError(w, err)
			return
		}

		// TODO: unique constraint handle username or email
		err = a.usersStorage.CreateUser(replacePasswordOnHash(&userData, passHash))
		if err != nil {
			web.WriteServerError(w, err)
			return
		}

		web.WriteJSON(w, struct{}{})
		return
	}

	web.WriteServerError(w, err)
}

func (a *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	panic("implement me")
}

func replacePasswordOnHash(user *models.RegistrationData, hash string) *models.RegistrationData {
	return &models.RegistrationData{
		Username: user.Username,
		Email:    user.Email,
		Password: hash,
	}
}
