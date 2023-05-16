package middleware

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
			w.WriteHeader(401)
			return
		}

		jwtStr := c.Value
		jwtOk := parseJWT(jwtStr)

		if !jwtOk {
			w.WriteHeader(401)
			return
		}

		// continue with the passing of the methods
		f(w, req)
	}
}

// parseJWT checks if the JWT is valid and in date
func parseJWT(jwtStr string) bool {
	tkn, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("SecretYouShouldHide"), nil
	})

	if err != nil {
		return false
	}

	if tkn.Method != jwt.SigningMethodHS256 {
		return false
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	inDate := jwtInDate(claims)
	if !ok || !inDate {
		return false
	}

	return true
}

// jwtInDate checks whether the JWT (claims) section is in date
func jwtInDate(claims jwt.MapClaims) bool {
	if int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
		return false
	}
	return true
}
