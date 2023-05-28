package auth

import (
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/user"
	blazem_user "blazem/pkg/domain/user"
	"blazem/pkg/domain/users"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/auth"

func TestAuth(t *testing.T) {

	type testcase struct {
		desc           string
		username       string
		password       string
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()

	responderMock := responder.NewResponder()
	usersMgrMock := users.NewUserStoreMock()
	jwtMgrMock := jwt_manager.NewJWTManager([]byte("SecretYouShouldHide"))

	usersMgrMock.Insert("user1", &user.User{
		Id:       "user1",
		Name:     "Jack",
		Role:     "admin",
		Username: "jack",
		Password: "password",
	})

	testcases := []testcase{
		{
			desc:           "HAPPY user that exists and valid",
			username:       "jack",
			password:       "password",
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			authMgrMock := NewAuthMgr(
				rMock,
				responderMock,
				usersMgrMock,
				jwtMgrMock,
			)
			authMgrMock.Register()

			w := httptest.NewRecorder()
			type authVal struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			body, _ := json.Marshal(authVal{
				Username: testCase.username,
				Password: testCase.password,
			})

			r, _ := http.NewRequest("POST", testRoute, bytes.NewBuffer(body))
			authMgrMock.Router.ServeHTTP(w, r)

			type userResponse struct {
				JWT  string            `json:"jwt"`
				User *blazem_user.User `json:"user"`
			}

			type authResponse struct {
				Msg  string `json:"msg"`
				Data userResponse
			}

			var response authResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
		})
	}
}
