package auth

import (
	"encoding/json"
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
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		web.WriteServerError(w, err)
		return
	}

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

	isExist, err := a.tokensStorage.CheckToken(candidate.ID)
	if isExist {
		err = a.tokensStorage.ChangeToken(candidate.ID, tkns.RefreshToken)
		if err != nil {
			web.WriteServerError(w, err)
			return
		}
	} else {
		err = a.tokensStorage.Create(candidate.ID, tkns.RefreshToken)
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

// Logout - curl -i -X POST -H "Authorization: Bearer <token>" http://localhost:5050/api/logout
func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	headerToken := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerToken, " ")

	if len(bearerToken) < 2 {
		web.WriteBadRequest(w, fmt.Errorf("bearer token not found"))
		return
	}

	userData, err := a.tokensManager.VerifyJWTToken(bearerToken[1])
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
	op := "services.auth.Registration"
	defer func() {
		err := r.Body.Close()
		if err != nil {
			a.logger.Warn(fmt.Sprintf("request body close error: %s", err.Error()), slog.String("op", op))
			return
		}
	}()

	var userData models.RegistrationData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	passHash, err := a.passManager.HashPassword(userData.Password)
	if err != nil {
		web.WriteServerError(w, err)
		return
	}

	err = a.usersStorage.CreateUser(replacePasswordOnHash(&userData, passHash))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			web.WriteBadRequest(w, fmt.Errorf("email must be unique"))
			return

		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			web.WriteBadRequest(w, fmt.Errorf("username must be unique"))
			return

		}
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, struct{}{})
	return
}

// Refresh - curl -i -X POST -H "Cookie: refresh_token=eyJhbG.eyJlwMD.Iuu4C4n" http://localhost:5050/api/refresh
func (a *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil {
		web.WriteBadRequest(w, fmt.Errorf("refresh token in cookey not found"))
		return
	}

	userData, err := a.tokensManager.VerifyJWTToken(c.Value)
	if err != nil {
		web.WriteBadRequest(w, fmt.Errorf("invalid token"))
		return
	}

	tkns, err := a.tokensManager.GenerateJWTTokens(userData.UserID, userData.UserRoleID)
	if err != nil {
		web.WriteServerError(w, err)
		return
	}

	err = a.tokensStorage.ChangeToken(userData.UserID, tkns.RefreshToken)
	if err != nil {
		web.WriteServerError(w, err)
		return
	}

	web.WriteJSON(w, &tokens.Tokens{
		AcceptToken:  tkns.AcceptToken,
		RefreshToken: tkns.RefreshToken,
	})
}

func replacePasswordOnHash(user *models.RegistrationData, hash string) *models.RegistrationData {
	return &models.RegistrationData{
		Username: user.Username,
		Email:    user.Email,
		Password: hash,
	}
}
