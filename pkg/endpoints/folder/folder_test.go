package folder

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"blazem/pkg/domain/user"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/folder/"

func TestFolder(t *testing.T) {

	type testcase struct {
		desc           string
		folderId       string
		expectedResult types.DataInFolder
		expectedStatus int
	}

	user := &user.User{
		Id:   "testuser",
		Name: "Jack",
		Role: "Admin",
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
	folderMap, _ := folder.FolderToMap(*folderMock)
	nodeMock.Data.Store("testdoc", folderMap)
	nodeMock.Data.Store("testdoc2", map[string]interface{}{
		"type":   "text",
		"key":    "testdoc2",
		"folder": "testdoc",
		"value":  "hello",
		"date":   "2023-05-27T22:37:30.025153+01:00",
	})

	responderMock := responder.NewResponder()
	dataStoreMock := storer.NewStore(nodeMock)
	jwtMgrMock := jwt_manager.NewJWTManager([]byte("SecretYouShouldHide"))

	testcases := []testcase{
		{
			desc:     "HAPPY provided correct folder id",
			folderId: "testdoc",
			expectedResult: types.DataInFolder{
				Data: []types.SendData{
					{
						Data: map[string]interface{}{
							"key":    "testdoc2",
							"date":   "2023-05-27T22:37:30.025153+01:00",
							"folder": "testdoc",
							"type":   "text",
							"value":  "hello",
						},
						Key: "testdoc2",
					},
				},
				FolderName: "testfolder",
			},
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			folderMgrMock := NewFolderMgr(
				rMock,
				nodeMock,
				responderMock,
				dataStoreMock,
				jwtMgrMock,
			)
			folderMgrMock.Register()

			w := httptest.NewRecorder()
			expiration := time.Now().Add(10 * time.Second)
			jwt, _ := jwtMgrMock.CreateJWT(user, expiration)

			cookie := http.Cookie{
				Name:    "token",
				Value:   jwt,
				Expires: expiration,
				Path:    "/",
			}

			r, _ := http.NewRequest("GET", testRoute+testCase.folderId, nil)
			r.AddCookie(&cookie)
			folderMgrMock.Router.ServeHTTP(w, r)

			type folderResponse struct {
				Msg  string             `json:"msg"`
				Data types.DataInFolder `json:"data"`
			}

			var response folderResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.DeepEqual(t, response.Data, testCase.expectedResult)
		})
	}
}
