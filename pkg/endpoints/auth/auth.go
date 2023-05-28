package auth

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/responder"
	blazem_user "blazem/pkg/domain/user"
	"blazem/pkg/domain/users"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var route = "/auth"

type AuthMgr struct {
	Router    *mux.Router
	Responder responder.Responder
	UserStore users.UserStorer
	JWTMgr    jwt_manager.JWTManager
}

func NewAuthMgr(router *mux.Router, responder responder.Responder, userStore users.UserStorer, jwtMgr jwt_manager.JWTManager) *AuthMgr {
	return &AuthMgr{
		Router:    router,
		Responder: responder,
		UserStore: userStore,
		JWTMgr:    jwtMgr,
	}
}

// Auth endpoint returns a JWT set for an expiration if the user exists
// it also sets a cookie for the client of this JWT
func (e *AuthMgr) Auth() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var authVal struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		json.NewDecoder(req.Body).Decode(&authVal)
		auth, err := authUser(e, authVal.Username, authVal.Password)
		if !auth {
			e.Responder.Error(w, 401, err)
			return
		}

		user, err := e.UserStore.GetByUsername(authVal.Username)
		if err != nil {
			e.Responder.Error(w, 401, err)
			return
		}

		expirationDate := time.Now().Add(10 * 24 * 60 * time.Minute)
		jwt, err := e.JWTMgr.CreateJWT(user, expirationDate)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   jwt,
			Expires: expirationDate,
		})

		type response struct {
			JWT  string            `json:"jwt"`
			User *blazem_user.User `json:"user"`
		}

		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg: "Successfully authenticated user",
			Data: response{
				JWT:  jwt,
				User: user,
			},
		})
	}
}

// authUser returns true if user is authed, false if not
func authUser(e *AuthMgr, username string, password string) (bool, error) {
	user, err := e.UserStore.GetByUsername(username)
	if err != nil {
		logger.Logger.Warn(err.Error())
		return false, err
	}
	if user.Password != password {
		err := errors.New("Incorrect username or password")
		logger.Logger.Warn(err.Error())
		return false, err
	}
	return true, nil
}

func (e *AuthMgr) Register() {
	e.Router.HandleFunc(route, e.Auth()).Methods("POST")
}
