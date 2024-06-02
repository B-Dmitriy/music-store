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
	"github.com/go-playground/validator/v10"

	tokensStore "github.com/B-Dmitriy/music-store/internal/storage/tokens"
)

type AuthService struct {
	logger        *slog.Logger
	passManager   *password.PasswordManager
	tokensManager *tokens.TokensManager
	usersStorage  *users.UsersStorage
	tokensStorage *tokensStore.TokenStorage
	validator     *validator.Validate
}

func New(
	logger *slog.Logger,
	pm *password.PasswordManager,
	tm *tokens.TokensManager,
	us *users.UsersStorage,
	ts *tokensStore.TokenStorage,
	v *validator.Validate,
) *AuthService {
	return &AuthService{
		logger:        logger,
		passManager:   pm,
		tokensManager: tm,
		usersStorage:  us,
		tokensStorage: ts,
		validator:     v,
	}
}

func (h *AuthService) panicRecover(w http.ResponseWriter, op string) {
	if r := recover(); r != nil {
		h.logger.Error("panic in services.auth", slog.String("op", op))
		web.WriteServerError(w, fmt.Errorf("server error"))
		return
	}
}

func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	op := "services.auth.Login"

	defer a.panicRecover(w, op)
	defer r.Body.Close()

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

func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	op := "services.auth.Logout"

	defer a.panicRecover(w, op)

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

func (a *AuthService) Registration(w http.ResponseWriter, r *http.Request) {
	op := "services.auth.Registration"

	defer a.panicRecover(w, op)
	defer r.Body.Close()

	var userData models.RegistrationData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	err = a.validator.Struct(&userData)
	if err != nil {
		web.WriteBadRequest(w, err.(validator.ValidationErrors))
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

func (a *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	op := "services.auth.Refresh"

	defer a.panicRecover(w, op)

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
