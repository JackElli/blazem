package endpoints

import (
	"blazem/pkg/domain/cors"
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	blazem_store "blazem/pkg/domain/storer"
	blazem_users "blazem/pkg/domain/users"
	"blazem/pkg/endpoints/adddoc"
	"blazem/pkg/endpoints/addfolder"
	"blazem/pkg/endpoints/adduser"
	"blazem/pkg/endpoints/auth"
	"blazem/pkg/endpoints/deletedoc"
	"blazem/pkg/endpoints/folder"
	"blazem/pkg/endpoints/folders"
	"blazem/pkg/endpoints/getdoc"
	"blazem/pkg/endpoints/getuser"
	"blazem/pkg/endpoints/middleware"
	"blazem/pkg/endpoints/nodemap"
	"blazem/pkg/endpoints/parent"
	"blazem/pkg/endpoints/permissions"
	"blazem/pkg/endpoints/query"
	"blazem/pkg/endpoints/recentquery"
	"blazem/pkg/endpoints/stats"
	"blazem/pkg/endpoints/users"
	blazem_query "blazem/pkg/query"
	"net/http"

	"github.com/gorilla/mux"
)

// Create all of the endpoints for Blazem
func SetupEndpoints(node *node.Node) error {

	responder := responder.NewResponder()
	dataStore := blazem_store.NewStore(node)
	jwtMgr := jwt_manager.NewJWTManager([]byte("SecretYouShouldHide"))
	queryer := blazem_query.NewQuery(node)
	userStore := blazem_users.NewUserStore()
	err := userStore.SetupUsers()
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	http.Handle("/", cors.CORS(r))

	middlewareMgr := middleware.NewMiddlewareMgr(jwtMgr)
	permissionsMgr := permissions.NewPermissionsMgr(userStore, jwtMgr)

	public := r.PathPrefix("/").Subrouter()
	authMgr := auth.NewAuthMgr(public, responder, userStore, jwtMgr)
	authMgr.Register()

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middlewareMgr.Middleware)
	folderMgr := folder.NewFolderMgr(protected, node, responder, dataStore, jwtMgr)
	folderMgr.Register()
	addFolderMgr := addfolder.NewAddFolderMgr(protected, node, responder, dataStore, jwtMgr)
	addFolderMgr.Register()
	getDocMgr := getdoc.NewGetDocMgr(protected, node, responder)
	getDocMgr.Register()
	addDocMgr := adddoc.NewAddDocMgr(protected, node, responder, dataStore)
	addDocMgr.Register()
	parentMgr := parent.NewParentMgr(protected, node, responder)
	parentMgr.Register()
	nodemapMgr := nodemap.NewNodemapMgr(protected, node, responder)
	nodemapMgr.Register()
	foldersMgr := folders.NewFoldersMgr(protected, node, responder)
	foldersMgr.Register()
	statsMgr := stats.NewStatsMgr(protected, node, responder)
	statsMgr.Register()
	queryMgr := query.NewQueryMgr(protected, node, responder, queryer)
	queryMgr.Register()
	recentQueryMgr := recentquery.NewRecentQueryMgr(protected, node, responder)
	recentQueryMgr.Register()
	getUsersMgr := users.NewUsersMgr(protected, responder, userStore)
	getUsersMgr.Register()
	getUserMgr := getuser.NewGetUserMgr(protected, responder, userStore)
	getUserMgr.Register()

	admin := r.PathPrefix("/").Subrouter()
	admin.Use(middlewareMgr.Middleware)
	admin.Use(permissionsMgr.Permissions)
	deleteDocMgr := deletedoc.NewDeleteDocMgr(admin, responder, dataStore)
	deleteDocMgr.Register()
	addUserMgr := adduser.NewAddUserMgr(admin, responder, userStore)
	addUserMgr.Register()

	return nil
}
