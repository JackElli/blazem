package deletedoc

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteDocMgr struct {
	Router    *mux.Router
	Responder responder.Responder
	DataStore storer.Storer
}

func NewDeleteDocMgr(router *mux.Router, responder responder.Responder, dataStore storer.Storer) *DeleteDocMgr {
	return &DeleteDocMgr{
		Router:    router,
		Responder: responder,
		DataStore: dataStore,
	}
}

// We want to delete a document from Blazem
func (e *DeleteDocMgr) DeleteDoc() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		docId := mux.Vars(req)["id"]
		docData, err := e.DataStore.Load(docId)
		if err != nil {
			e.Responder.Error(w, 404, errors.New("Cannot delete document, as it's not found"))
			return
		}

		doc := docData.(map[string]interface{})
		docKey := doc["key"].(string)
		folder := doc["folder"]

		err = e.DataStore.Delete(docKey, folder)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		global.DataChanged = true
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg: "Successfully deleted doc",
		})
	}
}

func (e *DeleteDocMgr) Register() {
	e.Router.HandleFunc("/doc/{id:[a-zA-Z0-9-]+}", e.DeleteDoc()).Methods("DELETE")
}
