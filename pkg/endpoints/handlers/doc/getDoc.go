package doc

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// We want to fetch a document and return it to the user
func GetDoc(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank == global.FOLLOWER {
			e.Responder.Error(w, 500, errors.New("Cannot fetch doc from a follower node"))
			return
		}

		docId := mux.Vars(req)["id"]
		getData, ok := e.Node.Data.Load(docId)
		if !ok {
			e.Responder.Error(w, 404, errors.New("Doc not found"))
			return
		}

		sendDataJson := formatData(getData.(map[string]interface{}), docId)
		sendData := endpoint.SendData{
			Key:  docId,
			Data: sendDataJson,
		}

		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg:  "Successfully retrieved doc",
			Data: sendData,
		})
	}
}

// We want to do a bit of manipulation to the document for instance
// we want to read from disk if it's not a text file or a folder
func formatData(getData map[string]interface{}, dataKey string) map[string]interface{} {
	if getData["type"] == "text" {
		return getData
	}

	if getData["type"] == "folder" {
		return getData
	}

	data, _ := ioutil.ReadFile("data/" + dataKey)
	var dataJSON global.JsonData
	json.Unmarshal(data, &dataJSON)
	getData["value"] = dataJSON["value"].(string)
	return getData
}
