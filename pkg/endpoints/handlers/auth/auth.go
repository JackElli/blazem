package auth

import (
	"blazem/pkg/domain/endpoint"
	"time"

	"net/http"

	"github.com/golang-jwt/jwt"
)

// Auth endpoint returns a JWT set for an expiration if the user exists
// it also sets a cookie for the client of this JWT
func Auth(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		auth := true
		if !auth {
			r.Respond(w, endpoint.EndpointResponse{
				Code: 401,
				Msg:  "User not authorised to do that",
			})
			return
		}

		secretkey := []byte("SecretYouShouldHide")
		expirationDate := time.Now().Add(20 * time.Minute)
		t := jwt.New(jwt.SigningMethodHS256)

		claims := t.Claims.(jwt.MapClaims)
		claims["exp"] = expirationDate.Unix()

		jwt, err := t.SignedString(secretkey)
		if err != nil {
			r.Respond(w, endpoint.EndpointResponse{
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

		r.Respond(w, endpoint.EndpointResponse{
			Code: 200,
			Msg:  "Successfully authenticated user",
			Data: jwt,
		})
	}
}
