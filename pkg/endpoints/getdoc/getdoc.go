package getdoc

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type GetDocMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
}

func NewGetDocMgr(router *mux.Router, node *node.Node, responder responder.Responder) *GetDocMgr {
	return &GetDocMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
	}
}

// We want to fetch a document and return it to the user
func (e *GetDocMgr) GetDoc() func(w http.ResponseWriter, req *http.Request) {
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

func (e *GetDocMgr) Register() {
	e.Router.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", e.GetDoc()).Methods("GET")
}
