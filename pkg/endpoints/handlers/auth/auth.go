package auth

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/logger"
	blazem_user "blazem/pkg/domain/user"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// Auth endpoint returns a JWT set for an expiration if the user exists
// it also sets a cookie for the client of this JWT
func Auth(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var authVal struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		json.NewDecoder(req.Body).Decode(&authVal)
		auth, err := authUser(e, authVal.Username, authVal.Password)

		if !auth {
			e.Responder.Respond(w, endpoint.EndpointResponse{
				Code: 401,
				Msg:  "User not authorised to do that, error: " + err.Error(),
			})
			return
		}

		user, err := e.UserStore.GetByUsername(authVal.Username)
		if err != nil {
			e.Responder.Respond(w, endpoint.EndpointResponse{
				Code: 401,
				Msg:  "User not authorised to do that, error: " + err.Error(),
			})
			return
		}

		secretkey := []byte("SecretYouShouldHide")
		expirationDate := time.Now().Add(10 * 24 * 60 * time.Minute)
		jwt, err := createJWT(secretkey, user, expirationDate)

		if err != nil {
			e.Responder.Respond(w, endpoint.EndpointResponse{
				Code: 500,
				Msg:  "Cannot auth user as jwt cannot be created as " + err.Error(),
			})
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

		e.Responder.Respond(w, endpoint.EndpointResponse{
			Code: 200,
			Msg:  "Successfully authenticated user",
			Data: response{
				JWT:  jwt,
				User: user,
			},
		})
	}
}

// authUser returns true if user is authed, false if not
func authUser(e *endpoint_manager.EndpointManager, username string, password string) (bool, error) {
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

// createJWT creates a JWT and returns the token and an error if there is one
func createJWT(secretkey []byte, user *blazem_user.User, expirationDate time.Time) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	claims := t.Claims.(jwt.MapClaims)
	claims["exp"] = expirationDate.Unix()
	claims["user"] = user.Id

	jwt, err := t.SignedString(secretkey)
	return jwt, err
}
