package auth

import (
	"fmt"
	"github.com/B-Dmitriy/music-store/pgk/web"
	"net/http"
	"strings"
)

func (a *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")
		bearerArr := strings.Split(headerToken, " ")

		if len(bearerArr) < 2 {
			web.WriteUnauthorized(w, fmt.Errorf("bearer token not found"))
			return
		}

		bearerToken := bearerArr[1]

		_, err := a.tokensManager.VerifyJWTToken(bearerToken)
		if err != nil {
			web.WriteUnauthorized(w, fmt.Errorf("token invalid"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
