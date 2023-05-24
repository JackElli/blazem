package doc

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// We want to delete a document from Blazem
func DeleteDoc(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		docId := mux.Vars(req)["id"]
		docData, docFound := e.Node.Data.Load(docId)
		if !docFound {
			e.Responder.Error(w, 404, errors.New("Cannot delete document, as it's not found"))
			return
		}

		doc := docData.(map[string]interface{})
		docKey := doc["key"].(string)
		folder := doc["folder"]

		err := e.DataStore.Delete(docKey, folder)
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
