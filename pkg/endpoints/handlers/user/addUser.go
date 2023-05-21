package user

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/user"
	"encoding/json"
	"net/http"
)

func AddUser(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var userToAdd user.User
		err := json.NewDecoder(req.Body).Decode(&userToAdd)
		if err != nil {
			e.Responder.Respond(w, endpoint.EndpointResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		err = e.Node.UserStore.Insert(userToAdd.Id, &userToAdd)
		if err != nil {
			e.Responder.Respond(w, endpoint.EndpointResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		e.Responder.Respond(w, endpoint.EndpointResponse{
			Code: 200,
			Msg:  "Successfully added user",
		})
	}
}
