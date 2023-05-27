package addfolder

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/folder"
	blazem_folder "blazem/pkg/domain/folder"
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"errors"

	"blazem/pkg/domain/global"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type AddFolderMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
	DataStore storer.Storer
	JWTMgr    jwt_manager.JWTManager
}

func NewAddFolderMgr(router *mux.Router, node *node.Node, responder responder.Responder, dataStore storer.Storer, jwtMgr jwt_manager.JWTManager) *AddFolderMgr {
	return &AddFolderMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
		DataStore: dataStore,
		JWTMgr:    jwtMgr,
	}
}

// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
// coming in, write to disk and add to the map
func (e *AddFolderMgr) AddFolder() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank != global.MASTER {
			e.Responder.Error(w, 500, errors.New("Should be master"))
			return
		}

		c, err := req.Cookie("token")
		if err != nil {
			e.Responder.Error(w, 500, errors.New("No user logged in"))
			return
		}

		var folder folder.Folder
		err = json.NewDecoder(req.Body).Decode(&folder)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		err = validate(folder)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		jwtVal := c.Value
		userId, err := e.JWTMgr.GetCurrentUserId(jwtVal)
		folder.DateCreated = time.Now().Format("2006-01-02T15:04:05")
		folder.CreatedBy = userId

		folderMap, err := blazem_folder.FolderToMap(folder)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		err = e.DataStore.Store(folder.Key, folder.Folder, folderMap)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		global.DataChanged = true
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg: "Added folder successfully",
		})
	}
}

// validate checks whether the incoming folder is a valid
// folder structure
func validate(folder folder.Folder) error {
	if folder.Name == "" ||
		folder.Key == "" {
		return errors.New("Folder is invalid")
	}

	return nil
}

func (e *AddFolderMgr) Register() {
	e.Router.HandleFunc("/folder", e.AddFolder()).Methods("POST")
}
