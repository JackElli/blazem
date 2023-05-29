package addfolder

import (
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"blazem/pkg/domain/user"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

var testRoute = "/folder"

func TestAddFolder(t *testing.T) {

	type testcase struct {
		desc           string
		expectedMsg    string
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()
	nodeMock.Rank = global.MASTER
	nodeMock.Data = sync.Map{}

	folderMock := folder.NewFolder(
		"testfolder",
		"Amazing Folder",
		true,
		"jack",
	)

	user := &user.User{
		Id:   "jack",
		Name: "Jack",
		Role: "Admin",
	}

	responderMock := responder.NewResponder()
	dataStoreMock := storer.NewStoreMock(nodeMock)
	jwtMgrMock := jwt_manager.NewJWTManager([]byte("SecretYouShouldHide"))

	expiration := time.Now().Add(10 * time.Second)
	jwt, _ := jwtMgrMock.CreateJWT(user, expiration)

	cookie := http.Cookie{
		Name:    "token",
		Value:   jwt,
		Expires: expiration,
		Path:    "/",
	}

	testcases := []testcase{
		{
			desc:           "HAPPY added correct folder",
			expectedMsg:    "Added folder successfully",
			expectedStatus: 200,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			addFolderMgrMock := NewAddFolderMgr(
				rMock,
				nodeMock,
				responderMock,
				dataStoreMock,
				jwtMgrMock,
			)
			addFolderMgrMock.Register()

			w := httptest.NewRecorder()

			body, _ := json.Marshal(folderMock)
			r, _ := http.NewRequest("POST", testRoute, bytes.NewBuffer(body))
			r.AddCookie(&cookie)
			addFolderMgrMock.Router.ServeHTTP(w, r)

			type addFolderResponse struct {
				Msg string `json:"msg"`
			}

			var response addFolderResponse
			json.NewDecoder(w.Body).Decode(&response)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
			assert.Equal(t, response.Msg, testCase.expectedMsg)
		})
	}
}
