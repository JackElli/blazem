package getuser

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/users"
	"net/http"

	"github.com/gorilla/mux"
)

var route = "/user/{id:[a-zA-Z0-9-:]+}"

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

		user, err := e.UserStore.Get(userId)
		if err != nil {
			e.Responder.Error(w, 404, err)
			return
		}

		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg:  "Successfully retrieved user",
			Data: user,
		})
	}
}

func (e *GetUserMgr) Register() {
	e.Router.HandleFunc(route, e.GetUser()).Methods("GET")
}
