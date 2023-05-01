package parent

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"net/http"
)

func NewParentHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return ParentHandler
}

func ParentHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &ParentEndpoint{
		Endpoint: *e,
	}
	return pe.getParentFolders
}

// Return the results of the nodemap to the client
func (e *ParentEndpoint) getParentFolders(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"all"})

	if req.Method != "GET" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	folderId := req.URL.Query().Get("folder")
	if folderId == "" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "No folder passed",
		})
		return
	}
	parents := GetParentFolders(e.Endpoint.Node, folderId)
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Parent folders retrieved Successfully",
		Data: parents,
	})
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
