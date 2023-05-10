package parent

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"net/http"

	"github.com/gorilla/mux"
)

func Parent(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		folderId := mux.Vars(req)["id"]
		if folderId == "" {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No folder passed",
			})
			return
		}
		parents := GetParentFolders(r.Node, folderId)
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Parent folders retrieved Successfully",
			Data: parents,
		})
	}
}

// This function returns all of the folders that parent the folder we are
// searching for recursively
func GetParentFolders(node *global.Node, searchFolderId string) []types.Folder {
	folderId := searchFolderId
	folders := make([]types.Folder, 0)
	for folderId != "" {
		folderInfo, ok := node.Data.Load(folderId)
		if !ok {
			folderId = ""
			continue
		}
		var folderMap = folderInfo.(global.Document)
		if folderId != searchFolderId {
			folders = append(folders, types.Folder{
				Folder:     "N/A",
				Key:        folderMap["key"].(string),
				FolderName: folderMap["folderName"].(string),
				DocCount:   -1,
			})
		}
		if folderMap["folder"] == nil {
			folderId = ""
			continue
		}
		folderId = folderMap["folder"].(string)
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
