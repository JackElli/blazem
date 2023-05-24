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
			e.Responder.Error(w, 500, err)
			return
		}

		err = e.UserStore.Insert(userToAdd.Id, &userToAdd)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg: "Successfully added user",
		})
	}
}
