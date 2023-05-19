package folder

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"errors"
	"log"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

// We want to return all of the data currently stored within this folder, including
// folders and data
func GetDataFolder(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		folderId := mux.Vars(req)["id"]
		if folderId == "" {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No folder passed",
			})
			return
		}
		folderName, err := GetFolderName(e.Node, folderId)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		var DOCS_TO_RENDER float64
		var returnData types.DataInFolder
		MAX_DOCS := 30
		dataInFolder := make([]types.SendData, 0)

		// Could do the sorting on the fly?
		// like a tree of some sort?
		e.Node.Data.Range(func(key, value interface{}) bool {
			doc := value.(global.Document)
			if _, ok := doc["folder"]; ok {
				if doc["folder"].(string) == folderId {
					docKey := doc["key"].(string)
					dataInFolder = append(dataInFolder, types.SendData{
						Key:  docKey,
						Data: doc,
					})
				}
			}
			return true
		})

		DOCS_TO_RENDER = math.Min(float64(MAX_DOCS), float64(len(dataInFolder)))

		sort.Slice(dataInFolder, func(i, j int) bool {
			if _, convOk := dataInFolder[i].Data["date"].(time.Time); !convOk {
				dateI, errI := time.Parse("2006-01-02T15:04:05", dataInFolder[i].Data["date"].(string))
				dateJ, errJ := time.Parse("2006-01-02T15:04:05", dataInFolder[j].Data["date"].(string))
				if errI != nil || errJ != nil {
					log.Fatal(errI)
				}
				return dateI.Unix() > dateJ.Unix()
			}
			return dataInFolder[i].Data["date"].(time.Time).Unix() > dataInFolder[j].Data["date"].(time.Time).Unix()
		})

		returnData.Data = dataInFolder[0:int(DOCS_TO_RENDER)]
		returnData.FolderName = folderName

		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved data in folder",
			Data: returnData,
		})
	}
}

// Returns the name of the folder, given the folderId
func GetFolderName(node *node.Node, folderId string) (string, error) {
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
