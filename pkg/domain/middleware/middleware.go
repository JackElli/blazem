package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// Middleware allows a http handler function to be passed through if a
// jwt is valid, if not, return a 401
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
		h.ServeHTTP(w, req)
	})
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

func GetCurrentUserId(jwtStr string) (string, error) {
	tkn, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("SecretYouShouldHide"), nil
	})
	if err != nil {
		return "", err
	}

	if tkn.Method != jwt.SigningMethodHS256 {
		return "", errors.New("not correct signing algorithm")
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("JWT not correct")
	}

	userId := claims["user"].(string)
	return userId, nil
}

// jwtInDate checks whether the JWT (claims) section is in date
func jwtInDate(claims jwt.MapClaims) bool {
	if int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
		return false
	}
	return true
}
