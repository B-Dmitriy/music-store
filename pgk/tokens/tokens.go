package tokens

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AcceptToken  string `json:"acceptToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserData struct {
	UserID     int `json:"userID"`
	UserRoleID int `json:"userRoleID"`
}

type TokensManager struct {
	secretKey []byte
}

func New(secretKey string) *TokensManager {
	return &TokensManager{
		secretKey: []byte(secretKey),
	}
}

func (tm *TokensManager) GenerateJWTTokens(userID int, userRoleID int) (*Tokens, error) {
	payloadAcceptToken := jwt.MapClaims{
		"userID":     userID,
		"userRoleID": userRoleID,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}

	payloadRefreshToken := jwt.MapClaims{
		"userID":     userID,
		"userRoleID": userRoleID,
		"exp":        time.Now().Add(time.Hour * 24 * 14).Unix(),
	}

	acceptToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payloadAcceptToken)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payloadRefreshToken)

	at, err := acceptToken.SignedString(tm.secretKey)
	if err != nil {
		return nil, err
	}

	rt, err := refreshToken.SignedString(tm.secretKey)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AcceptToken:  at,
		RefreshToken: rt,
	}, nil
}

func (tm *TokensManager) VerifyJWTToken(tokenString string) (*UserData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, err := strconv.Atoi(fmt.Sprintf("%v", claims["userID"]))
		if err != nil {
			return nil, err
		}
		userRoleID, err := strconv.Atoi(fmt.Sprintf("%v", claims["userRoleID"]))
		if err != nil {
			return nil, err
		}
		return &UserData{
			UserRoleID: userRoleID,
			UserID:     userID,
		}, nil
	} else {
		return nil, err
	}
}
