package query

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/query"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/query"

func TestQuery(t *testing.T) {

	type testcase struct {
		desc              string
		queryStr          string
		expectedResultLen int
		expectedStatus    int
	}

	nodeMock := node.NewNode()
	nodeMock.Rank = global.MASTER
	nodeMock.SetupLogger()

	nodeMock.Data.Store("testdoc", map[string]interface{}{
		"key":   "testdoc",
		"type":  "text",
		"value": "Hello",
	})

	nodeMock.Data.Store("testdoc2", map[string]interface{}{
		"key":   "testdoc2",
		"type":  "text",
		"value": "test",
	})

	responderMock := responder.NewResponder()

	queryMock := query.NewQuery(nodeMock)

	rMock := mux.NewRouter()
	getQueryMock := NewQueryMgr(
		rMock,
		nodeMock,
		responderMock,
		queryMock,
	)
	getQueryMock.Register()

	testcases := []testcase{
		{
			desc:              "HAPPY select all",
			queryStr:          "SELECT all",
			expectedResultLen: 2,
			expectedStatus:    200,
		},
		{
			desc:              "HAPPY like query",
			queryStr:          "SELECT all WHERE value LIKE \"test\"",
			expectedResultLen: 1,
			expectedStatus:    200,
		},
		{
			desc:              "NEGATIVE incorrect query",
			queryStr:          "SELEC",
			expectedResultLen: 0,
			expectedStatus:    500,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {

			w := httptest.NewRecorder()

			type queryVal struct {
				Query string `json:"query"`
			}

			body, _ := json.Marshal(queryVal{
				Query: testCase.queryStr,
			})

			r, _ := http.NewRequest("POST", testRoute, bytes.NewBuffer(body))

			getQueryMock.Router.ServeHTTP(w, r)

			type queryResponse struct {
				Msg  string              `json:"msg"`
				Data types.SendQueryData `json:"data"`
			}

			var response queryResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, len(response.Data.Docs), testCase.expectedResultLen)
		})
	}
}
