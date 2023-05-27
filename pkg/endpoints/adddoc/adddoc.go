package adddoc

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type AddDocMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
	DataStore storer.Storer
}

func NewAddDocMgr(router *mux.Router, node *node.Node, responder responder.Responder, dataStore storer.Storer) *AddDocMgr {
	return &AddDocMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
		DataStore: dataStore,
	}
}

// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
// coming in, write to disk and add to the map
func (e *AddDocMgr) AddDoc() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank != global.MASTER {
			e.Responder.Error(w, 500, errors.New("Should be master"))
			return
		}
		var doc map[string]interface{}
		json.NewDecoder(req.Body).Decode(&doc)

		err := validateDoc(doc)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		document := transformNewDoc(e.Node, doc)
		docKey := doc["key"].(string)
		folder := doc["folder"]

		e.Node.WriteDocToDisk(document)
		err = e.DataStore.Store(docKey, folder, document)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		global.DataChanged = true
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg: "Added document successfully",
		})
	}
}

// validate checks whether the doc is valid and able to
// be added to blazem
func validateDoc(doc map[string]interface{}) error {
	_, hasKey := doc["key"]

	if !hasKey {
		return errors.New("This doc has no key")
	}

	return nil
}

// transformNewDoc transforms the document coming in, to something that is optimised and
// info-full for retrieval
func transformNewDoc(node *node.Node,
	dataToAdd map[string]interface{}) map[string]interface{} {
	document := dataToAdd

	document["date"] = time.Now().Format("2006-01-02T15:04:05")

	_, typeExist := document["type"]
	if !typeExist {
		// should read file type
		document["type"] = "text"
	}

	loadDoc, updateDoc := node.Data.Load(dataToAdd["key"])
	if updateDoc {
		document["date"] = loadDoc.(map[string]interface{})["date"].(string)
	}

	return document
}

func (e *AddDocMgr) Register() {
	e.Router.HandleFunc("/doc", e.AddDoc()).Methods("POST")
}
