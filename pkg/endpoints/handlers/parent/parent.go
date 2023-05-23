package parent

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/node"
	"net/http"

	"github.com/gorilla/mux"
)

func Parent(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		folderId := mux.Vars(req)["id"]
		parents := GetParentFolders(e.Node, folderId)
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Parent folders retrieved Successfully",
			Data: parents,
		})
	}
}

// This function returns all of the folders that parent the folder we are
// searching for recursively
func GetParentFolders(node *node.Node, searchFolderId string) []types.Folder {
	folderId := searchFolderId
	folders := make([]types.Folder, 0)
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
func reverse(lst []types.Folder) []types.Folder {
	newLst := make([]types.Folder, 0)
	for i := len(lst) - 1; i >= 0; i-- {
		newLst = append(newLst, lst[i])
	}
	return newLst
}
