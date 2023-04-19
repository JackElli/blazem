package handlers

import (
	"blazem/global"
	"net/http"
)

func FolderHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.folderHandler
}

// We want to return all of the root folders in the data i.e every folder
// that doesnt have a folder parent. We fetch the folder names, add them to the
// folder map and add the corresponding global.Document count
func (node *Node) folderHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, nil)

	if req.Method != "GET" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	var folders = node.getAllFolders()
	folders = node.getFolderDocCount(folders)
	folders = node.storeDocCount(folders)

	JsonResponse(w, EndpointResponse{
		200,
		"Successfully retrieved folders",
		folders,
	})
}

// We want to get all of the folders currently in Blazem
func (node *Node) getAllFolders() map[string]Folder {
	var folders = make(map[string]Folder, 0)

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

		folders[folderKey] = Folder{
			inFolder,
			folderKey,
			folderName,
			0,
			backedUp,
		}
		return true
	})
	return folders
}

// We want to get all of the folder doc counts
func (node *Node) getFolderDocCount(folders map[string]Folder) map[string]Folder {
	node.Data.Range(func(k, value interface{}) bool {
		if folder, exists := value.(global.Document)["folder"].(string); exists && folder != "" {
			currDocCount := folders[folder].DocCount
			folders[folder] = Folder{
				folders[folder].Folder,
				folders[folder].Key,
				folders[folder].FolderName,
				currDocCount + 1,
				folders[folder].BackedUp,
			}
		}
		return true
	})
	return folders
}

// We want to store the doc counts in blazem
func (node *Node) storeDocCount(folders map[string]Folder) map[string]Folder {
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
