package handlers

import (
	"blazem/global"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"
)

func GetDataInFolder(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.getDataInFolderHandler
}

func (node *Node) getDataInFolderHandler(w http.ResponseWriter, req *http.Request) {
	// We want to return all of the data currently stored within this folder, including
	// folders and data

	// Please can we fix this function
	WriteHeaders(w, nil)

	var returnData DataInFolder
	var folderId = req.URL.Query().Get("folder")
	var user = req.URL.Query().Get("user")
	var folderName = node.GetFolderName(folderId)
	var nodeData = make([]global.Document, lenOfSyncMap(node.Data))
	var dataInFolder []SendData = []SendData{}
	var dataInd = 0

	if user != "jack" {
		json.NewEncoder(w).Encode("no auth")
		return
	}

	node.Data.Range(func(key, value interface{}) bool {
		nodeData[dataInd] = value.(global.Document)
		dataInd++
		return true
	})

	sort.Slice(nodeData, func(i, j int) bool {
		if _, convOk := nodeData[i]["date"].(time.Time); !convOk {
			dateI, errI := time.Parse("2006-01-02T15:04:05", nodeData[i]["date"].(string))
			dateJ, errJ := time.Parse("2006-01-02T15:04:05", nodeData[j]["date"].(string))
			if errI != nil || errJ != nil {
				log.Fatal(errI)
			}
			return dateI.Unix() > dateJ.Unix()
		}
		return nodeData[i]["date"].(time.Time).Unix() > nodeData[j]["date"].(time.Time).Unix()
	})

	numOfItems := 0
	for i, data := range nodeData {
		key := nodeData[i]["key"].(string)
		if data["type"] != "text" {
			data["value"] = "file"
		}
		if numOfItems == 40 {
			break
		}
		if data["folder"] == folderId {
			sendData := SendData{key, data}
			dataInFolder = append(dataInFolder, sendData)
			numOfItems++
		}
	}

	sort.Slice(dataInFolder[:], func(i, j int) bool {
		iRunes := []rune(dataInFolder[i].Data["type"].(string))
		jRunes := []rune(dataInFolder[j].Data["type"].(string))
		iVal := int(iRunes[0])
		jVal := int(jRunes[0])
		return iVal < jVal
	})

	returnData.Data = dataInFolder
	returnData.FolderName = folderName
	returnData.ParentFolders = node.getParentFolders(folderId)

	json.NewEncoder(w).Encode(returnData)
}

func (node *Node) GetFolderName(folderId string) string {
	folder, ok := node.Data.Load(folderId)
	if !ok {
		return ""
	}
	folderMap := folder.(global.Document)
	return folderMap["folderName"].(string)
}

func (node *Node) getParentFolders(searchFolderId string) []Folder {
	// This function returns all of the folders that parent the folder we are
	// searching for recursively
	var folderId = searchFolderId
	var folders []Folder = []Folder{}
	for folderId != "" {
		folderInfo, ok := node.Data.Load(folderId)
		if !ok {
			folderId = ""
			continue
		}
		folderMap := folderInfo.(global.Document)
		if folderId != searchFolderId {
			folders = append(folders, Folder{
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

func reverse(lst []Folder) []Folder {
	// Reverse order of list
	var newLst []Folder = []Folder{}
	for i := len(lst) - 1; i >= 0; i-- {
		newLst = append(newLst, lst[i])
	}
	return newLst
}
