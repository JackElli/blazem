package folder

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"net/http"
)

func NewFolderHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return FolderHandler
}

func FolderHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	ce := &FolderEndpoint{
		Endpoint: *e,
	}
	return ce.folderHandler
}

// We want to return all of the root folders in the data i.e every folder
// that doesnt have a folder parent. We fetch the folder names, add them to the
// folder map and add the corresponding global.Document count
func (e *FolderEndpoint) folderHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, nil)

	if req.Method != "GET" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	folders := GetAllFolders(e.Endpoint.Node)
	folders = GetFolderDocCount(e.Endpoint.Node, folders)
	folders = StoreDocCount(e.Endpoint.Node, folders)
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully retrieved folders",
		Data: folders,
	})
}

// We want to get all of the folders currently in Blazem
func GetAllFolders(node *global.Node) map[string]types.Folder {
	folders := make(map[string]types.Folder, 0)
	node.Data.Range(func(k, value interface{}) bool {
		dataType := value.(global.Document)["type"]
		if dataType != "folder" {
			return true
		}
		folderKey := value.(global.Document)["key"].(string)
		folderName := value.(global.Document)["folderName"].(string)

		var inFolder string
		var exists bool
		var backedUp bool = false
		if value.(global.Document)["backedUp"] != nil {
			backedUp = value.(global.Document)["backedUp"].(bool)
		}
		if inFolder, exists = value.(global.Document)["folder"].(string); !exists {
			inFolder = ""
		}
		folders[folderKey] = types.Folder{
			Folder:     inFolder,
			Key:        folderKey,
			FolderName: folderName,
			DocCount:   0,
			BackedUp:   backedUp,
		}
		return true
	})
	return folders
}

// We want to get all of the folder doc counts
func GetFolderDocCount(node *global.Node, folders map[string]types.Folder) map[string]types.Folder {
	node.Data.Range(func(k, value interface{}) bool {
		if folder, exists := value.(global.Document)["folder"].(string); exists && folder != "" {
			currDocCount := folders[folder].DocCount
			folders[folder] = types.Folder{
				Folder:     folders[folder].Folder,
				Key:        folders[folder].Key,
				FolderName: folders[folder].FolderName,
				DocCount:   currDocCount + 1,
				BackedUp:   folders[folder].BackedUp,
			}
		}
		return true
	})
	return folders
}

// We want to store the doc counts in blazem
func StoreDocCount(node *global.Node, folders map[string]types.Folder) map[string]types.Folder {
	for _, folder := range folders {
		if folder.Folder == "" {
			continue
		}
		folderData, _ := node.Data.Load(folder.Key)
		folderData.(global.Document)["docCount"] = folder.DocCount
		node.Data.Store(folder.Key, folderData)
		delete(folders, folder.Key)
	}
	return folders
}