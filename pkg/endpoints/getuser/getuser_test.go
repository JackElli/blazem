package getuser

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
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

var testRoute = "/user/"

func TestConnect(t *testing.T) {

	type testcase struct {
		desc           string
		userId         string
		expectedResult user.User
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()
	nodeMock.Rank = global.MASTER

	responderMock := responder.NewResponder()

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
	testcases := []testcase{
		{
			desc:   "HAPPY retrieved user that exists",
			userId: "user1",
			expectedResult: user.User{
				Id:       "user1",
				Name:     "Jack",
				Role:     "admin",
				Username: "jacktest",
				Password: "Supersecure",
			},
			expectedStatus: 200,
		},
		{
			desc:           "NEGATIVE retrieved user that doesn't exist",
			userId:         "testing",
			expectedStatus: 404,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			getUserMgrMock := NewGetUserMgr(
				rMock,
				responderMock,
				userStoreMock,
			)
			getUserMgrMock.Register()

			w := httptest.NewRecorder()

			r, _ := http.NewRequest("GET", testRoute+testCase.userId, nil)
			getUserMgrMock.Router.ServeHTTP(w, r)

			type getUserResponse struct {
				Msg  string    `json:"msg"`
				Data user.User `json:"data"`
			}

			var response getUserResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.Equal(t, response.Data, testCase.expectedResult)
		})
	}
}
