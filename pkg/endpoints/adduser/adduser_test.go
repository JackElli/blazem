package adduser

import (
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/user"
	"blazem/pkg/domain/users"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/user"

func TestUsers(t *testing.T) {

	type testcase struct {
		desc           string
		expectedMsg    string
		expectedStatus int
	}

	userStoreMock := users.NewUserStoreMock()
	userStoreMock.Users = map[string]user.User{}
	responderMock := responder.NewResponder()

	rMock := mux.NewRouter()

	usersMgrMock := NewAddUserMgr(
		rMock,
		responderMock,
		userStoreMock,
	)
	usersMgrMock.Register()

	user := &user.User{
		Id:       "jack",
		Name:     "Jack",
		Role:     "Admin",
		Username: "jack",
		Password: "testpass",
	}

	testcases := []testcase{
		{
			desc:           "HAPPY successful addition of a user",
			expectedMsg:    "Successfully added user",
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {

			w := httptest.NewRecorder()
			body, _ := json.Marshal(user)
			r, _ := http.NewRequest("POST", testRoute, bytes.NewBuffer(body))

			usersMgrMock.Router.ServeHTTP(w, r)

			type addUserResponse struct {
				Msg string `json:"msg"`
			}

			var response addUserResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.Equal(t, response.Msg, testCase.expectedMsg)
		})
	}
}
