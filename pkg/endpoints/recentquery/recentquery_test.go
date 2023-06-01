package recentquery

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

var testRoute = "/recentQueries"

func TestConnect(t *testing.T) {

	type testcase struct {
		desc           string
		recentQueries  map[string]string
		expectedResult map[string]string
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()
	nodeMock.Rank = global.MASTER

	responderMock := responder.NewResponder()

	testcases := []testcase{
		{
			desc: "HAPPY retrieved one recent query",
			recentQueries: map[string]string{
				"test": "SELECT",
			},
			expectedResult: map[string]string{
				"test": "SELECT",
			},
			expectedStatus: 200,
		},
		{
			desc: "HAPPY retrieved one recent queries",
			recentQueries: map[string]string{
				"test":  "SELECT all",
				"asd":   "SELECT all WHERE jack",
				"213":   "SELECT 2",
				"asddd": "SELECT 4",
			},
			expectedResult: map[string]string{
				"test":  "SELECT all",
				"asd":   "SELECT all WHERE jack",
				"213":   "SELECT 2",
				"asddd": "SELECT 4",
			},
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			nodeMock.RecentQueries = testCase.recentQueries

			getRecentQueriesMgrMock := NewRecentQueryMgr(
				rMock,
				nodeMock,
				responderMock,
			)
			getRecentQueriesMgrMock.Register()

			w := httptest.NewRecorder()

			r, _ := http.NewRequest("GET", testRoute, nil)
			getRecentQueriesMgrMock.Router.ServeHTTP(w, r)

			type getRecentQueriesResponse struct {
				Msg  string            `json:"msg"`
				Data map[string]string `json:"data"`
			}

			var response getRecentQueriesResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, response.Data, testCase.expectedResult)
		})
	}
}
