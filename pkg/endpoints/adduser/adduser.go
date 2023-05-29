package adduser

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/user"
	"blazem/pkg/domain/users"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var route = "/user"

type AddUserMgr struct {
	Router    *mux.Router
	Responder responder.Responder
	UserStore users.UserStorer
}

func NewAddUserMgr(router *mux.Router, responder responder.Responder, userStore users.UserStorer) *AddUserMgr {
	return &AddUserMgr{
		Router:    router,
		Responder: responder,
		UserStore: userStore,
	}
}

func (e *AddUserMgr) AddUser() func(w http.ResponseWriter, req *http.Request) {
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

func (e *AddUserMgr) Register() {
	e.Router.HandleFunc(route, e.AddUser()).Methods("POST")
}
