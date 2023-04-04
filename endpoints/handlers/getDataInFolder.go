package handlers

import (
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

	// Please can we fix this function

	WriteHeaders(w, nil)

	var returnData DataInFolder

	folderId := req.URL.Query().Get("folder")
	user := req.URL.Query().Get("user")
	folderName := node.GetFolderName(folderId)

	if user != "jack" {
		json.NewEncoder(w).Encode("no auth")
		return
	}

	// need to sort data by date
	// breaking change, as added new JSON field
	nodeData := make([]map[string]interface{}, lenOfSyncMap(node.Data))

	dataInd := 0
	node.Data.Range(func(key, value interface{}) bool {
		nodeData[dataInd] = value.(map[string]interface{})
		dataInd++
		return true
	})

	// not sure why this is like this?
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

	var dataInFolder []SendData = []SendData{}
	numOfItems := 0
	for i, data := range nodeData {
		key := nodeData[i]["key"].(string)

		// push the word file instead
		// of file data so we save
		// on loading times
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
	folderMap := folder.(map[string]interface{})
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
		folderMap := folderInfo.(map[string]interface{})
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
	var newLst []Folder = []Folder{}
	for i := len(lst) - 1; i >= 0; i-- {
		newLst = append(newLst, lst[i])
	}
	return newLst
}
