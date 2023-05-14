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
		jwt, err := createJWT(secretkey, expirationDate)

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

// createJWT creates a JWT and returns the token and an error if there is one
func createJWT(secretkey []byte, expirationDate time.Time) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	claims := t.Claims.(jwt.MapClaims)
	claims["exp"] = expirationDate.Unix()

	jwt, err := t.SignedString(secretkey)
	return jwt, err
}
