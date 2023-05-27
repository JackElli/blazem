package middleware

import (
	"blazem/pkg/domain/jwt_manager"
	"net/http"
)

type MiddlewareMgr struct {
	JWTMgr jwt_manager.JWTManager
}

func NewMiddlewareMgr(jwtMgr jwt_manager.JWTManager) *MiddlewareMgr {
	return &MiddlewareMgr{
		JWTMgr: jwtMgr,
	}
}

// Middleware allows a http handler function to be passed through if a
// jwt is valid, if not, return a 401
func (e *MiddlewareMgr) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("token")
		if err != nil {
			w.WriteHeader(401)
			return
		}

		jwtStr := c.Value
		_, jwtOk := e.JWTMgr.ParseJWT(jwtStr)

		if !jwtOk {
			w.WriteHeader(401)
			return
		}

		// continue with the passing of the methods
		h.ServeHTTP(w, req)
	})
}
