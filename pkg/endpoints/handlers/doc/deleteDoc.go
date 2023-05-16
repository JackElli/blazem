package doc

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"net/http"

	"github.com/gorilla/mux"
)

// We want to delete a document from Blazem
func DeleteDoc(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		docId := mux.Vars(req)["id"]
		if docId == "" {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Doc key not provided",
			})
			return
		}
		_, docFound := e.Node.Data.Load(docId)
		if !docFound {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 404,
				Msg:  "Cannot delete document, as it's not found",
			})
			return
		}
		e.Node.Data.Delete(docId)
		global.DataChanged = true
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully deleted doc",
		})
	}
}
