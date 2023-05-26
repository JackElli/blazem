package endpoint_manager

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/node"
	"net/http"
)

// We want to return all of the root folders in the data i.e every folder
// that doesnt have a folder parent. We fetch the folder names, add them to the
// folder map and add the corresponding document count
func (e *EndpointManager) Folders() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		folders := GetAllFolders(e.Node)
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg:  "Successfully retrieved folders",
			Data: folders,
		})
	}
}

// We want to get all of the folders currently in Blazem
func GetAllFolders(node *node.Node) map[string]types.Folder {
	folders := make(map[string]types.Folder, 0)
	node.Data.Range(func(k, value interface{}) bool {
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
