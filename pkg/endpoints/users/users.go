package users

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/users"
	"net/http"

	"github.com/gorilla/mux"
)

type UsersMgr struct {
	Router    *mux.Router
	Responder responder.Responder
	UserStore users.UserStorer
}

func NewUsersMgr(router *mux.Router, responder responder.Responder, userStore users.UserStorer) *UsersMgr {
	return &UsersMgr{
		Router:    router,
		Responder: responder,
		UserStore: userStore,
	}
}

func (e *UsersMgr) GetUsers() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg:  "Successfully retrieved users",
			Data: e.UserStore.List(),
		})
	}
}

func (e *UsersMgr) Register() {
	e.Router.HandleFunc("/users", e.GetUsers()).Methods("GET")
}
