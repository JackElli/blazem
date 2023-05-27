package getuser

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/users"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type GetUserMgr struct {
	Router    *mux.Router
	Responder responder.Responder
	UserStore users.UserStorer
}

func NewGetUserMgr(router *mux.Router, responder responder.Responder, userStore users.UserStorer) *GetUserMgr {
	return &GetUserMgr{
		Router:    router,
		Responder: responder,
		UserStore: userStore,
	}
}

func (e *GetUserMgr) GetUser() func(w http.ResponseWriter, req *http.Request) {
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

func (e *GetUserMgr) Register() {
	e.Router.HandleFunc("/user/{id:[a-zA-Z0-9-:]+}", e.GetUser()).Methods("GET")
}
