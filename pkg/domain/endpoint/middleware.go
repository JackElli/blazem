package endpoint

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// Middleware allows a http handler function to be passed through if a
// jwt is valid, if not, return a 401
func Middleware(f func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("token")
		if err != nil {
			w.WriteHeader(403)
			return
		}

		jwtStr := c.Value
		tkn, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
			return t, nil
		})

		if tkn.Method != jwt.SigningMethodHS256 {
			w.WriteHeader(403)
			return
		}

		claims, ok := tkn.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(403)
			return
		}

		// check if token is in date
		if int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
			w.WriteHeader(403)
			return
		}

		// continue with the passing of the methods
		f(w, req)
	}
}
