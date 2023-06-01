package connect

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/connect/"

func TestConnect(t *testing.T) {

	type testcase struct {
		desc           string
		ip             string
		expectedMsg    string
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()
	nodeMock.Rank = global.MASTER

	responderMock := responder.NewResponder()
	testcases := []testcase{
		{
			desc:           "HAPPY connected successfully",
			ip:             "192.168.1.something",
			expectedMsg:    "Successfully connected",
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			addDocMgrMock := NewConnectMgr(
				rMock,
				nodeMock,
				responderMock,
			)
			addDocMgrMock.Register()

			w := httptest.NewRecorder()

			r, _ := http.NewRequest("POST", testRoute+testCase.ip, nil)
			addDocMgrMock.Router.ServeHTTP(w, r)

			type connectResponse struct {
				Msg string `json:"msg"`
			}

			var response connectResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.Equal(t, response.Msg, testCase.expectedMsg)
		})
	}
}
