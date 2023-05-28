package getdoc

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/doc/"

func TestGetDoc(t *testing.T) {

	type testcase struct {
		desc           string
		docId          string
		expectedResult interface{}
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.Rank = global.MASTER
	nodeMock.Data = sync.Map{}
	nodeMock.Data.Store("testdoc", map[string]interface{}{
		"type":  "text",
		"key":   "testdoc",
		"value": "Hello",
	})

	responderMock := responder.NewResponder()

	rMock := mux.NewRouter()

	getDocMgrMock := NewGetDocMgr(
		rMock,
		nodeMock,
		responderMock,
	)
	getDocMgrMock.Register()

	testcases := []testcase{
		{
			desc:  "HAPPY provided correct doc id",
			docId: "testdoc",
			expectedResult: endpoint.SendData{
				Data: map[string]interface{}{
					"type":  "text",
					"key":   "testdoc",
					"value": "Hello",
				},
				Key: "testdoc",
			},
			expectedStatus: 200,
		},
		{
			desc:  "NEGATIVE provided incorrect doc id",
			docId: "testdoc2",
			expectedResult: endpoint.SendData{
				Data: nil,
			},
			expectedStatus: 404,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", testRoute+testCase.docId, nil)

			getDocMgrMock.Router.ServeHTTP(w, r)

			type docResponse struct {
				Msg  string            `json:"msg"`
				Data endpoint.SendData `json:"data"`
			}

			var response docResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, response.Data, testCase.expectedResult)
		})
	}
}
