package middlewares

import (
	"fooder/api/utils"
	"fooder/config"
	"fooder/repositories/models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func NewAuthMiddleware(c *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenHeader := r.Header.Get("Authorization")

			if tokenHeader == "" {
				errResponse := utils.NewApiError(http.StatusUnauthorized, "Missing auth token")
				utils.SendError(errResponse, w)
				return
			}

			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				errResponse := utils.NewApiError(http.StatusBadRequest, "Missing auth token")
				utils.SendError(errResponse, w)
				return
			}

			tokenPart := splitted[1]
			tk := &models.Token{}

			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(c.Auth.JWTPassword), nil
			})

			if err != nil {
				errResponse := utils.NewApiError(http.StatusUnauthorized, "Malformed authentication token")
				utils.SendError(errResponse, w)
				return
			}

			if !token.Valid {
				errResponse := utils.NewApiError(http.StatusForbidden, "Token is invalid")
				utils.SendError(errResponse, w)
				return
			}

			r = r.WithContext(utils.CreateContextWithToken(r.Context(), tk))
			next.ServeHTTP(w, r)
		})
	}
}
