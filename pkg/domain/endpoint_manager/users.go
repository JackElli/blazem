package endpoint_manager

import (
	"blazem/pkg/domain/endpoint"
	"net/http"
)

func (e *EndpointManager) GetUsers() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg:  "Successfully retrieved users",
			Data: e.UserStore.List(),
		})
	}
}
