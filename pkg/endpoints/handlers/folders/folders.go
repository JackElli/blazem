package folders

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/node"
	"net/http"
)

// We want to return all of the root folders in the data i.e every folder
// that doesnt have a folder parent. We fetch the folder names, add them to the
// folder map and add the corresponding document count
func Folders(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		folders := GetAllFolders(e.Node)
		folders = GetFolderDocCount(e.Node, folders)
		folders = StoreDocCount(e.Node, folders)
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved folders",
			Data: folders,
		})
	}
}

// We want to get all of the folders currently in Blazem
func GetAllFolders(node *node.Node) map[string]types.Folder {
	folders := make(map[string]types.Folder, 0)
	node.Data.Range(func(k, value interface{}) bool {
		dataType := value.(map[string]interface{})["type"]
		if dataType != "folder" {
			return true
		}
		folderKey := value.(map[string]interface{})["key"].(string)
		folderName := value.(map[string]interface{})["folderName"].(string)

		var inFolder string
		var exists bool
		var backedUp bool = false
		if value.(map[string]interface{})["backedUp"] != nil {
			backedUp = value.(map[string]interface{})["backedUp"].(bool)
		}
		if inFolder, exists = value.(map[string]interface{})["folder"].(string); !exists {
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
func GetFolderDocCount(node *node.Node, folders map[string]types.Folder) map[string]types.Folder {
	node.Data.Range(func(k, value interface{}) bool {
		if folder, exists := value.(map[string]interface{})["folder"].(string); exists && folder != "" {
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
func StoreDocCount(node *node.Node, folders map[string]types.Folder) map[string]types.Folder {
	for _, folder := range folders {
		if folder.Folder == "" {
			continue
		}
		folderData, _ := node.Data.Load(folder.Key)
		folderData.(map[string]interface{})["docCount"] = folder.DocCount
		node.Data.Store(folder.Key, folderData)
		delete(folders, folder.Key)
	}
	return folders
}
