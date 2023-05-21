package endpoints

import (
	"blazem/pkg/domain/cors"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/middleware"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
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
	"blazem/pkg/endpoints/handlers/user"
	"blazem/pkg/endpoints/handlers/users"
	blazem_query "blazem/pkg/query"
	"net/http"

	"github.com/gorilla/mux"
)

// Create all of the endpoints for Blazem
func SetupEndpoints(node *node.Node) error {
	endpointMgr := endpoint_manager.NewEndpointManager(
		node,
		responder.NewResponder(),
		blazem_query.NewQuery(node),
		storer.NewStore(node),
	)

	r := mux.NewRouter()
	http.Handle("/", cors.CORS(r))

	public := r.PathPrefix("/").Subrouter()
	public.HandleFunc("/auth", auth.Auth(endpointMgr)).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.Middleware)
	protected.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", doc.GetDoc(endpointMgr)).Methods("GET")
	protected.HandleFunc("/doc", doc.AddDoc(endpointMgr)).Methods("POST")
	protected.HandleFunc("/folder/{id:[a-zA-Z0-9-]+}", folder.GetFolderData(endpointMgr)).Methods("GET")
	protected.HandleFunc("/folder", folder.AddFolder(endpointMgr)).Methods("POST")
	protected.HandleFunc("/parents/{id:[a-zA-Z0-9-]+}", parent.Parent(endpointMgr)).Methods("GET")
	protected.HandleFunc("/nodemap", nodemap.NodeMap(endpointMgr)).Methods("GET")
	protected.HandleFunc("/folders", folders.Folders(endpointMgr)).Methods("GET")
	protected.HandleFunc("/stats", stats.Stats(endpointMgr)).Methods("GET")
	protected.HandleFunc("/query", query.Query(endpointMgr)).Methods("POST")
	protected.HandleFunc("/recentQueries", recentquery.RecentQuery(endpointMgr)).Methods("GET")
	protected.HandleFunc("/ping", ping.Ping(endpointMgr)).Methods("POST")
	protected.HandleFunc("/users", users.GetUsers(endpointMgr)).Methods("GET")
	protected.HandleFunc("/user/{id:[a-zA-Z0-9-:]+}", user.GetUser(endpointMgr)).Methods("GET")
	protected.HandleFunc("/user", user.AddUser(endpointMgr)).Methods("POST")

	admin := r.PathPrefix("/").Subrouter()
	admin.Use(middleware.Middleware)
	admin.Use(endpointMgr.Permissions)
	admin.HandleFunc("/connect/{ip:[a-zA-Z0-9-.-]+}", connect.Connect(endpointMgr)).Methods("POST")
	admin.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", doc.DeleteDoc(endpointMgr)).Methods("DELETE")
	admin.HandleFunc("/node/{ip:[a-zA-Z0-9-]+}", nodes.RemoveNode(endpointMgr)).Methods("DELETE")

	return nil
}
