package nodemap

import (
	"blazem/pkg/domain/endpoint"
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

var testRoute = "/nodemap"

func TestGetDoc(t *testing.T) {

	type testcase struct {
		desc           string
		expectedResult interface{}
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.Ip = "testip"
	nodeMock.Rank = global.MASTER
	nodeMock.NodeMap = []*node.Node{
		nodeMock,
	}

	responderMock := responder.NewResponder()

	rMock := mux.NewRouter()

	getNodemapMock := NewNodemapMgr(
		rMock,
		nodeMock,
		responderMock,
	)
	getNodemapMock.Register()

	testcases := []testcase{
		{
			desc: "HAPPY provided correct doc id",
			expectedResult: []endpoint.WebNodeMap{
				{
					Ip:     "testip",
					Active: true,
					Rank:   global.MASTER,
				},
			},
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", testRoute, nil)

			getNodemapMock.Router.ServeHTTP(w, r)

			type nodemapResponse struct {
				Msg  string                `json:"msg"`
				Data []endpoint.WebNodeMap `json:"data"`
			}

			var response nodemapResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, response.Data, testCase.expectedResult)
		})
	}
}
