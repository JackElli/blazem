package permissions

import "net/http"

var admin = false

func Permissions(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if admin {
			h.ServeHTTP(w, req)
			return
		}
		w.WriteHeader(403)
	})
}
