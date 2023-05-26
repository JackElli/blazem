package endpoint_manager

import (
	"blazem/pkg/domain/endpoint"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func (e *EndpointManager) GetUser() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userId := mux.Vars(req)["id"]
		if userId == "" {
			e.Responder.Error(w, 500, errors.New("Doc key not provided"))
			return
		}

		user, err := e.UserStore.Get(userId)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg:  "Successfully retrieved user",
			Data: user,
		})
	}
}
