package adddoc

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/doc"

func TestAddDoc(t *testing.T) {

	type testcase struct {
		desc           string
		expectedMsg    string
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()
	nodeMock.Rank = global.MASTER
	nodeMock.Data = sync.Map{}

	docMock := map[string]interface{}{
		"key":   "testdoc",
		"type":  "text",
		"value": "Testing adding doc",
	}

	responderMock := responder.NewResponder()
	dataStoreMock := storer.NewStoreMock(nodeMock)

	testcases := []testcase{
		{
			desc:           "HAPPY added correct document",
			expectedMsg:    "Added document successfully",
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			addDocMgrMock := NewAddDocMgr(
				rMock,
				nodeMock,
				responderMock,
				dataStoreMock,
			)
			addDocMgrMock.Register()

			w := httptest.NewRecorder()

			body, _ := json.Marshal(docMock)
			r, _ := http.NewRequest("POST", testRoute, bytes.NewBuffer(body))
			addDocMgrMock.Router.ServeHTTP(w, r)

			type addDocResponse struct {
				Msg string `json:"msg"`
			}

			var response addDocResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.Equal(t, response.Msg, testCase.expectedMsg)
		})
	}
}
