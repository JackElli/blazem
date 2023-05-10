package doc

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"net/http"
)

// We want to delete a document from Blazem
func DeleteDoc(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r.WriteHeaders(w, []string{"all"})
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if req.Method != "DELETE" {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Wrong method",
			})
			return
		}
		docKey := req.URL.Query().Get("key")
		if docKey == "" {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Doc key not provided",
			})
			return
		}
		_, docFound := r.Node.Data.Load(docKey)
		if !docFound {
			r.Respond(w, types.EndpointResponse{
				Code: 404,
				Msg:  "Cannot delete document, as it's not found",
			})
			return
		}
		r.Node.Data.Delete(docKey)
		global.DataChanged = true
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully deleted doc",
		})
	}
}
