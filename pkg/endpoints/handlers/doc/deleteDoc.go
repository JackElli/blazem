package doc

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/responder"
	"net/http"

	"github.com/gorilla/mux"
)

// We want to delete a document from Blazem
func DeleteDoc(r *responder.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		docId := mux.Vars(req)["id"]
		if docId == "" {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Doc key not provided",
			})
			return
		}
		_, docFound := r.Node.Data.Load(docId)
		if !docFound {
			r.Respond(w, types.EndpointResponse{
				Code: 404,
				Msg:  "Cannot delete document, as it's not found",
			})
			return
		}
		r.Node.Data.Delete(docId)
		global.DataChanged = true
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully deleted doc",
		})
	}
}
