package handlers

import (
	"blazem/global"
	"errors"
	"log"
	"net/http"
	"sort"
	"time"
)

func GetDataInFolder(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.getDataInFolderHandler
}

// We want to return all of the data currently stored within this folder, including
// folders and data
func (node *Node) getDataInFolderHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, nil)

	if req.Method != "GET" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	var returnData DataInFolder
	var folderId = req.URL.Query().Get("folder")
	if folderId == "" {
		JsonResponse(w, EndpointResponse{
			500,
			"No folder passed",
			nil,
		})
		return
	}
	var folderName, err = node.GetFolderName(folderId)
	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			err.Error(),
			nil,
		})
		return
	}
	var nodeData = make([]global.Document, lenOfSyncMap(node.Data))
	var dataInFolder = make([]SendData, 0)
	var dataInd = 0

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

	var numOfItems = 0
	for i, data := range nodeData {
		var key = nodeData[i]["key"].(string)
		if data["type"] != "text" {
			data["value"] = "file"
		}
		if numOfItems == 40 {
			break
		}
		if data["folder"] == folderId {
			var sendData = SendData{key, data}
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

	JsonResponse(w, EndpointResponse{
		200,
		"Successfully retrieved data in folder",
		returnData,
	})
}

// Returns the name of the folder, given the folderId
func (node *Node) GetFolderName(folderId string) (string, error) {
	folder, ok := node.Data.Load(folderId)
	if !ok {
		return "", errors.New("No document with that key")
	}
	var folderMap = folder.(global.Document)
	if folderMap["type"] != "folder" {
		return "", errors.New("No folder with that key")
	}
	return folderMap["folderName"].(string), nil
}

// This function returns all of the folders that parent the folder we are
// searching for recursively
func (node *Node) getParentFolders(searchFolderId string) []Folder {
	var folderId = searchFolderId
	var folders = make([]Folder, 0)
	for folderId != "" {
		var folderInfo, ok = node.Data.Load(folderId)
		if !ok {
			folderId = ""
			continue
		}
		var folderMap = folderInfo.(global.Document)
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
	var newLst = make([]Folder, 0)
	for i := len(lst) - 1; i >= 0; i-- {
		newLst = append(newLst, lst[i])
	}
	return newLst
}
