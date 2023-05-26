package endpoints

import (
	"blazem/pkg/domain/cors"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	blazem_store "blazem/pkg/domain/storer"
	blazem_users "blazem/pkg/domain/users"
	blazem_query "blazem/pkg/query"
	"net/http"

	"github.com/gorilla/mux"
)

// Create all of the endpoints for Blazem
func SetupEndpoints(node *node.Node) error {
	endpointMgr := endpoint_manager.NewEndpointManager(
		node,
		responder.NewResponder(),
		jwt_manager.NewJWTManager([]byte("SecretYouShouldHide")),
		blazem_query.NewQuery(node),
		blazem_users.NewUserStore(),
		blazem_store.NewStore(node),
	)

	err := endpointMgr.UserStore.SetupUsers()
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	http.Handle("/", cors.CORS(r))

	public := r.PathPrefix("/").Subrouter()
	public.HandleFunc("/auth", endpointMgr.Auth()).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(endpointMgr.Middleware)
	protected.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", endpointMgr.GetDoc()).Methods("GET")
	protected.HandleFunc("/doc", endpointMgr.AddDoc()).Methods("POST")
	protected.HandleFunc("/folder/{id:[a-zA-Z0-9-]+}", endpointMgr.GetFolderData()).Methods("GET")
	protected.HandleFunc("/folder", endpointMgr.AddFolder()).Methods("POST")
	protected.HandleFunc("/parents/{id:[a-zA-Z0-9-]+}", endpointMgr.Parent()).Methods("GET")
	protected.HandleFunc("/nodemap", endpointMgr.NodeMap()).Methods("GET")
	protected.HandleFunc("/folders", endpointMgr.Folders()).Methods("GET")
	protected.HandleFunc("/stats", endpointMgr.Stats()).Methods("GET")
	protected.HandleFunc("/query", endpointMgr.CompleteQuery()).Methods("POST")
	protected.HandleFunc("/recentQueries", endpointMgr.RecentQuery()).Methods("GET")
	protected.HandleFunc("/ping", endpointMgr.Ping()).Methods("POST")
	protected.HandleFunc("/users", endpointMgr.GetUsers()).Methods("GET")
	protected.HandleFunc("/user/{id:[a-zA-Z0-9-:]+}", endpointMgr.GetUser()).Methods("GET")

	admin := r.PathPrefix("/").Subrouter()
	admin.Use(endpointMgr.Middleware)
	admin.Use(endpointMgr.Permissions)
	admin.HandleFunc("/connect/{ip:[a-zA-Z0-9-.-]+}", endpointMgr.Connect()).Methods("POST")
	admin.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", endpointMgr.DeleteDoc()).Methods("DELETE")
	admin.HandleFunc("/node/{ip:[a-zA-Z0-9-]+}", endpointMgr.RemoveNode()).Methods("DELETE")
	admin.HandleFunc("/user", endpointMgr.AddUser()).Methods("POST")

	return nil
}
