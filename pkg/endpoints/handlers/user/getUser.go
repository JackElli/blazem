package user

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUser(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userId := mux.Vars(req)["id"]
		if userId == "" {
			e.Responder.Respond(w, endpoint.EndpointResponse{
				Code: 500,
				Msg:  "Doc key not provided",
			})
			return
		}

		user, err := e.Node.UserStore.Get(userId)
		if err != nil {
			e.Responder.Respond(w, endpoint.EndpointResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		e.Responder.Respond(w, endpoint.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved user",
			Data: user,
		})
	}
}
