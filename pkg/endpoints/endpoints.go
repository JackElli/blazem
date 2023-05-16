package endpoints

import (
	"blazem/pkg/domain/cors"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/middleware"
	"blazem/pkg/domain/responder"
	"blazem/pkg/endpoints/handlers/auth"
	"blazem/pkg/endpoints/handlers/connect"
	"blazem/pkg/endpoints/handlers/doc"
	"blazem/pkg/endpoints/handlers/folder"
	"blazem/pkg/endpoints/handlers/folders"
	"blazem/pkg/endpoints/handlers/nodemap"
	"blazem/pkg/endpoints/handlers/nodes"
	"blazem/pkg/endpoints/handlers/parent"
	"blazem/pkg/endpoints/handlers/ping"
	"blazem/pkg/endpoints/handlers/query"
	"blazem/pkg/endpoints/handlers/recentquery"
	"blazem/pkg/endpoints/handlers/stats"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Create all of the endpoints for Blazem
func SetupEndpoints(node *global.Node) error {
	responder := responder.NewResponder(node)
	r := mux.NewRouter()

	r.HandleFunc("/connect/{ip:[a-zA-Z0-9-.-]+}", middleware.Middleware(connect.Connect(responder))).Methods("POST")
	r.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", middleware.Middleware(doc.GetDoc(responder))).Methods("GET")
	r.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", middleware.Middleware(doc.DeleteDoc(responder))).Methods("DELETE")
	r.HandleFunc("/node/{ip:[a-zA-Z0-9-]+}", middleware.Middleware(nodes.RemoveNode(responder))).Methods("DELETE")
	r.HandleFunc("/folder/{id:[a-zA-Z0-9-]+}", middleware.Middleware(folder.GetDataFolder(responder))).Methods("GET")
	r.HandleFunc("/parents/{id:[a-zA-Z0-9-]+}", middleware.Middleware(parent.Parent(responder))).Methods("GET")
	r.HandleFunc("/nodemap", middleware.Middleware(nodemap.NodeMap(responder))).Methods("GET")
	r.HandleFunc("/doc", middleware.Middleware(doc.AddDoc(responder))).Methods("POST")
	r.HandleFunc("/folders", middleware.Middleware(folders.Folders(responder))).Methods("GET")
	r.HandleFunc("/stats", middleware.Middleware(stats.Stats(responder))).Methods("GET")
	r.HandleFunc("/query", middleware.Middleware(query.Query(responder))).Methods("POST")
	r.HandleFunc("/recentQueries", middleware.Middleware(recentquery.RecentQuery(responder))).Methods("GET")
	r.HandleFunc("/ping", middleware.Middleware(ping.Ping(responder))).Methods("POST")
	r.HandleFunc("/auth", auth.Auth(responder)).Methods("POST")

	setCorsMethods(r)
	return nil
}

// setCorsMethods allows us to choose which headers are allowed
func setCorsMethods(r *mux.Router) {
	allowedMethods := cors.GetAllowedMethods()
	allowedHeaders := cors.GetAllowedHeaders()
	allowedOrigins := cors.GetAllowedOrigin()
	exposeHeader := cors.GetExposedHeaders()

	http.Handle("/", handlers.CORS(
		handlers.AllowedMethods(allowedMethods),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders(exposeHeader),
		handlers.AllowedHeaders(allowedHeaders),
	)(r))
}
