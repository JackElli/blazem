package users

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"net/http"
)

func GetUsers(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		e.Responder.Respond(w, endpoint.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved users",
			Data: e.Node.UserStore.Users,
		})
	}
}
