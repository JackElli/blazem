package users

import (
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/user"
	"blazem/pkg/domain/users"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/users"

func TestUsers(t *testing.T) {

	type testcase struct {
		desc              string
		expectedResultLen int
		expectedStatus    int
	}

	userStoreMock := users.NewUserStore()
	userStoreMock.Users = map[string]user.User{
		"user1": {
			Id:       "user1",
			Name:     "Jack",
			Role:     "admin",
			Username: "jacktest",
			Password: "Supersecure",
		},
		"user2": {
			Id:       "user2",
			Name:     "Test",
			Role:     "basic",
			Username: "stest",
			Password: "Supersecure2",
		},
	}
	responderMock := responder.NewResponder()

	rMock := mux.NewRouter()

	usersMgrMock := NewUsersMgr(
		rMock,
		responderMock,
		userStoreMock,
	)
	usersMgrMock.Register()

	testcases := []testcase{
		{
			desc:              "HAPPY list users",
			expectedResultLen: 2,
			expectedStatus:    200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", testRoute, nil)

			usersMgrMock.Router.ServeHTTP(w, r)

			type usersResponse struct {
				Msg  string               `json:"msg"`
				Data map[string]user.User `json:"data"`
			}

			var response usersResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, len(response.Data), testCase.expectedResultLen)
		})
	}
}
