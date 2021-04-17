package middlewares

import (
	"fmt"
	"fooder/api/consts"
	"fooder/api/utils"
	"net/http"
)

func NewSessionMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(consts.SessionCookie)
			if err != nil {
				fmt.Printf("Error during session parsing: %s\n", err)
			}

			if cookie != nil {
				r = r.WithContext(utils.CreateContextWithSessionData(r.Context(), cookie.Value, r.Header.Get("Origin")))
			}
			next.ServeHTTP(w, r)
		})
	}
}
