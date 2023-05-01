package datainfolder

import (
	types "blazem/domain/endpoint"
	"blazem/domain/global"
	"errors"
	"log"
	"net/http"
	"sort"
	"time"
)

func NewGetDataFolderHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return GetDataFolderHandler
}

func GetDataFolderHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	de := &DataInFolderEndpoint{
		Endpoint: *e,
	}
	return de.getDataInFolderHandler
}

// We want to return all of the data currently stored within this folder, including
// folders and data
func (e *DataInFolderEndpoint) getDataInFolderHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, nil)
	if req.Method != "GET" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	var returnData types.DataInFolder
	folderId := req.URL.Query().Get("folder")
	if folderId == "" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "No folder passed",
		})
		return
	}
	folderName, err := GetFolderName(e.Endpoint.Node, folderId)
	if err != nil {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	nodeData := make([]global.Document, types.LenOfSyncMap(e.Endpoint.Node.Data))
	dataInFolder := make([]types.SendData, 0)
	dataInd := 0
	e.Endpoint.Node.Data.Range(func(key, value interface{}) bool {
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
		var key = nodeData[i]["key"].(string)
		if data["type"] != "text" {
			data["value"] = "file"
		}
		if numOfItems == 40 {
			break
		}
		if data["folder"] == folderId {
			sendData := types.SendData{
				Key:  key,
				Data: data,
			}
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

	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully retrieved data in folder",
		Data: returnData,
	})
}

// Returns the name of the folder, given the folderId
func GetFolderName(node *global.Node, folderId string) (string, error) {
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
