package cors

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func GetAllowedMethods() []string {
	return []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
}

func GetAllowedHeaders() []string {
	return []string{"Origin", "Content-Length", "Content-Type"}
}

func GetAllowedOrigin() []string {
	return []string{"http://localhost:5173"}
}

func GetExposedHeaders() []string {
	return []string{"Set-Cookie"}
}

// setCorsMethods allows us to choose which headers are allowed
func CORS(r *mux.Router) http.Handler {
	allowedMethods := GetAllowedMethods()
	allowedHeaders := GetAllowedHeaders()
	allowedOrigins := GetAllowedOrigin()
	exposeHeader := GetExposedHeaders()

	return handlers.CORS(
		handlers.AllowedMethods(allowedMethods),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders(exposeHeader),
		handlers.AllowedHeaders(allowedHeaders),
	)(r)
}
