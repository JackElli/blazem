package handlers

import (
	"blazem/global"
	"encoding/json"
	"net/http"
)

func FolderHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.folderHandler
}

func (node *Node) folderHandler(w http.ResponseWriter, req *http.Request) {
	// We want to return all of the root folders in the data i.e every folder
	// that doesnt have a folder parent. We fetch the folder names, add them to the
	// folder map and add the corresponding global.Document count
	WriteHeaders(w, nil)
	var folders = make(map[string]Folder, 0)

	// SPLIT THESE UP INTO SEPARATE FUNCS
	node.Data.Range(func(k, value interface{}) bool {
		dataType := value.(global.Document)["type"]
		if dataType == "folder" {
			var inFolder string
			var exists bool
			var backedUp bool = false
			folderKey := value.(global.Document)["key"].(string)
			folderName := value.(global.Document)["folderName"].(string)
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
		}
		return true
	})

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
	for _, folder := range folders {
		if folder.Folder != "" {
			// This stores doc count
			folderData, _ := node.Data.Load(folder.Key)
			folderData.(global.Document)["docCount"] = folder.DocCount
			node.Data.Store(folder.Key, folderData)
			delete(folders, folder.Key)
		}
	}
	json.NewEncoder(w).Encode(folders)
}
