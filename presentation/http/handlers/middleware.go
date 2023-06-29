package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
)

// BasicAuth middleware set up Basic Auth
func (rst *RestFulAPIs) BasicAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				username, password, ok := r.BasicAuth()
				if ok {
					usernameHash := sha256.Sum256([]byte(username))
					passwordHash := sha256.Sum256([]byte(password))
					expectedUserNameHash := sha256.Sum256(
						[]byte(rst.auth.username),
					)
					expectedPasswordHash := sha256.Sum256(
						[]byte(rst.auth.password),
					)
					usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:],
						expectedUserNameHash[:]) == 1)
					passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:],
						expectedPasswordHash[:]) == 1)

					if usernameMatch && passwordMatch {
						next.ServeHTTP(w, r)
						return
					}
				}

				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="utf-8"`)
				RespondWithError(
					w,
					http.StatusUnauthorized,
					fmt.Errorf("unauthorized"),
				)

			},
		)
	}
}
