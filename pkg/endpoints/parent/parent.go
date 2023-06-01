package parent

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"net/http"

	"github.com/gorilla/mux"
)

var route = "/parents/{id:[a-zA-Z0-9-]+}"

type ParentMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
}

func NewParentMgr(router *mux.Router, node *node.Node, responder responder.Responder) *ParentMgr {
	return &ParentMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
	}
}

func (e *ParentMgr) Parent() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		folderId := mux.Vars(req)["id"]
		parents := GetParentFolders(e.Node, folderId)
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg:  "Parent folders retrieved Successfully",
			Data: parents,
		})
	}
}

// This function returns all of the folders that parent the folder we are
// searching for recursively
func GetParentFolders(node *node.Node, searchFolderId string) []folder.Folder {
	folderId := searchFolderId
	folders := make([]folder.Folder, 0)
	for folderId != "" {
		doc, ok := node.Data.Load(folderId)
		if !ok {
			folderId = ""
			continue
		}
		folderData := doc.(map[string]interface{})
		folder, isFolder := folder.IsFolder(folderData)
		if !isFolder {
			folderId = ""
		}
		if folderId != searchFolderId {
			folders = append(folders, *folder)
		}
		if folder.Folder == "" {
			folderId = ""
			continue
		}
		folderId = folder.Folder
	}
	return reverse(folders)
}

// Reverse order of list
func reverse(lst []folder.Folder) []folder.Folder {
	newLst := make([]folder.Folder, 0)
	for i := len(lst) - 1; i >= 0; i-- {
		newLst = append(newLst, lst[i])
	}
	return newLst
}

func (e *ParentMgr) Register() {
	e.Router.HandleFunc(route, e.Parent()).Methods("GET")
}
