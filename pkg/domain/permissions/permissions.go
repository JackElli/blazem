package permissions

import (
	"blazem/pkg/domain/user"
	"net/http"
)

var CurrentUser *user.User

func Permissions(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if CurrentUser == nil {
			w.WriteHeader(401)
			return
		}
		if CurrentUser.Role == "admin" {
			h.ServeHTTP(w, req)
			return
		}
		w.WriteHeader(403)
	})
}
