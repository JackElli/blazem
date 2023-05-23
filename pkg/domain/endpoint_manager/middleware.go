package endpoint_manager

import (
	"errors"
	"net/http"
)

// Middleware allows a http handler function to be passed through if a
// jwt is valid, if not, return a 401
func (em *EndpointManager) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("token")
		if err != nil {
			w.WriteHeader(401)
			return
		}

		jwtStr := c.Value
		_, jwtOk := em.JWTManager.ParseJWT(jwtStr)

		if !jwtOk {
			w.WriteHeader(401)
			return
		}

		// continue with the passing of the methods
		h.ServeHTTP(w, req)
	})
}

// GetCurrentUserId returns the Id of the user who is logged in to
// this session
func (em *EndpointManager) GetCurrentUserId(jwtStr string) (string, error) {
	parseClaims, valid := em.JWTManager.ParseJWT(jwtStr)
	if !valid {
		return "", errors.New("JWT not valid")
	}

	claims := *parseClaims
	userId := claims["user"].(string)

	return userId, nil
}
