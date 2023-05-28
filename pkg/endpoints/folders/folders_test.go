package folders

import (
	"blazem/pkg/domain/folder"
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

var testRoute = "/folders"

func TestFolders(t *testing.T) {

	type testcase struct {
		desc              string
		expectedResultLen int
		expectedStatus    int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()
	nodeMock.Rank = global.MASTER
	nodeMock.Data = sync.Map{}

	folderMock := folder.NewFolder(
		"testdoc",
		"testfolder",
		true,
		"testuser",
	)
	folderMock2 := folder.NewFolder(
		"testdoc2",
		"testfolder2",
		true,
		"testuser",
	)
	folderMap, _ := folder.FolderToMap(*folderMock)
	folderMap2, _ := folder.FolderToMap(*folderMock2)

	nodeMock.Data.Store("testdoc", folderMap)
	nodeMock.Data.Store("testdoc2", folderMap2)

	responderMock := responder.NewResponder()

	testcases := []testcase{
		{
			desc:              "HAPPY get folders",
			expectedResultLen: 2,
			expectedStatus:    200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			folderMgrMock := NewFoldersMgr(
				rMock,
				nodeMock,
				responderMock,
			)
			folderMgrMock.Register()

			w := httptest.NewRecorder()

			r, _ := http.NewRequest("GET", testRoute, nil)
			folderMgrMock.Router.ServeHTTP(w, r)

			type foldersResponse struct {
				Msg  string                   `json:"msg"`
				Data map[string]folder.Folder `json:"data"`
			}

			var response foldersResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, len(response.Data), testCase.expectedResultLen)
		})
	}
}
