package parent

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

var testRoute = "/parents/"

func TestParent(t *testing.T) {

	type testcase struct {
		desc           string
		folderId       string
		expectedResult []folder.Folder
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()
	nodeMock.Data = sync.Map{}
	nodeMock.Rank = global.MASTER

	folderMock1 := folder.NewFolder(
		"testdoc",
		"testfolder",
		true,
		"testuser",
	)
	folderMock1.DateCreated = "2006-01-02T15:04:05"
	folderMock2 := folder.NewFolder(
		"testdoc2",
		"testfolder2",
		true,
		"testuser",
	)
	folderMock2.Folder = "testdoc"
	folderMap1, _ := folder.FolderToMap(*folderMock1)
	folderMap2, _ := folder.FolderToMap(*folderMock2)
	nodeMock.Data.Store("testdoc", folderMap1)
	nodeMock.Data.Store("testdoc2", folderMap2)

	responderMock := responder.NewResponder()

	rMock := mux.NewRouter()

	parentMock := NewParentMgr(
		rMock,
		nodeMock,
		responderMock,
	)
	parentMock.Register()

	testcases := []testcase{
		{
			desc:     "HAPPY provided correct doc id",
			folderId: "testdoc2",
			expectedResult: []folder.Folder{
				{
					Key:         "testdoc",
					Name:        "testfolder",
					CreatedBy:   "testuser",
					Global:      true,
					DateCreated: "2006-01-02T15:04:05",
					Type:        "folder",
				},
			},
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", testRoute+testCase.folderId, nil)

			parentMock.Router.ServeHTTP(w, r)

			type parentResponse struct {
				Msg  string          `json:"msg"`
				Data []folder.Folder `json:"data"`
			}

			var response parentResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, response.Data, testCase.expectedResult)
		})
	}
}
