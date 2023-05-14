package endpoints

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
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

type EndpointRouter struct {
	node *global.Node
	r    *mux.Router
}

func NewEndpointRouter(node *global.Node, router *mux.Router) *EndpointRouter {
	return &EndpointRouter{
		node: node,
		r:    router,
	}
}

// Create all of the endpoints for Blazem
func SetupEndpoints(node *global.Node) error {
	responder := endpoint.NewResponder(node)

	r := mux.NewRouter()
	r.HandleFunc("/connect/{ip:[a-zA-Z0-9-.-]+}", endpoint.Middleware(connect.Connect(responder))).Methods("POST")
	r.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", endpoint.Middleware(doc.GetDoc(responder))).Methods("GET")
	r.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", endpoint.Middleware(doc.DeleteDoc(responder))).Methods("DELETE")
	r.HandleFunc("/node/{ip:[a-zA-Z0-9-]+}", endpoint.Middleware(nodes.RemoveNode(responder))).Methods("DELETE")
	r.HandleFunc("/folder/{id:[a-zA-Z0-9-]+}", endpoint.Middleware(folder.GetDataFolder(responder))).Methods("GET")
	r.HandleFunc("/parents/{id:[a-zA-Z0-9-]+}", endpoint.Middleware(parent.Parent(responder))).Methods("GET")
	r.HandleFunc("/nodemap", endpoint.Middleware(nodemap.NodeMap(responder))).Methods("GET")
	r.HandleFunc("/doc", endpoint.Middleware(doc.AddDoc(responder))).Methods("POST")
	r.HandleFunc("/folders", endpoint.Middleware(folders.Folders(responder))).Methods("GET")
	r.HandleFunc("/stats", endpoint.Middleware(stats.Stats(responder))).Methods("GET")
	r.HandleFunc("/query", endpoint.Middleware(query.Query(responder))).Methods("POST")
	r.HandleFunc("/recentQueries", endpoint.Middleware(recentquery.RecentQuery(responder))).Methods("GET")
	r.HandleFunc("/ping", endpoint.Middleware(ping.Ping(responder))).Methods("POST")
	r.HandleFunc("/auth", auth.Auth(responder)).Methods("GET")

	allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "HEAD"}
	allowedHeaders := []string{"Origin", "Content-Length", "Content-Type"}
	allowedOrigins := []string{"http://localhost:5173"}
	http.Handle("/", handlers.CORS(
		handlers.AllowedMethods(allowedMethods),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders(allowedHeaders))(r),
	)
	return nil
}
