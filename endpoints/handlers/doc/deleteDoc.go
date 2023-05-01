package doc

import (
	types "blazem/domain/endpoint"
	"blazem/domain/global"
	"net/http"
)

func NewDeleteDocHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return DeleteDocHandler
}

func DeleteDocHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	de := &DocEndpoint{
		Endpoint: *e,
	}
	return de.deleteDocHandler
}

// We want to delete a document from Blazem
func (e *DocEndpoint) deleteDocHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{})

	if req.Method != "DELETE" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	docKey := req.URL.Query().Get("key")
	if docKey == "" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Doc key not provided",
		})
		return
	}
	_, docFound := e.Endpoint.Node.Data.Load(docKey)
	if !docFound {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 404,
			Msg:  "Cannot delete document, as it's not found",
		})
		return
	}
	e.Endpoint.Node.Data.Delete(docKey)
	global.DataChanged = true
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully deleted doc",
	})
}
