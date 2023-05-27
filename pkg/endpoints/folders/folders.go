package folders

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"net/http"

	"github.com/gorilla/mux"
)

type FoldersMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
}

func NewFoldersMgr(router *mux.Router, node *node.Node, responder responder.Responder) *FoldersMgr {
	return &FoldersMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
	}
}

// We want to return all of the root folders in the data i.e every folder
// that doesnt have a folder parent. We fetch the folder names, add them to the
// folder map and add the corresponding document count
func (e *FoldersMgr) Folders() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		folders := e.GetAllFolders()
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg:  "Successfully retrieved folders",
			Data: folders,
		})
	}
}

// We want to get all of the folders currently in Blazem
func (e *FoldersMgr) GetAllFolders() map[string]types.Folder {
	folders := make(map[string]types.Folder, 0)
	e.Node.Data.Range(func(k, value interface{}) bool {
		doc := value.(map[string]interface{})
		folder, isFolder := folder.IsFolder(doc)

		if !isFolder || folder.Folder != "" {
			return true
		}

		folders[folder.Key] = *folder

		return true
	})
	return folders
}

func (e *FoldersMgr) Register() {
	e.Router.HandleFunc("/folders", e.Folders()).Methods("GET")
}
