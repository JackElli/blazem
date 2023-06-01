package permissions

import (
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/user"
	"blazem/pkg/domain/users"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"gotest.tools/v3/assert"
)

func TestPermissions(t *testing.T) {
	type testcase struct {
		desc           string
		cookie         *http.Cookie
		expectedStatus int
	}

	nodeMock := node.NewNode()
	nodeMock.SetupLogger()

	userStoreMock := users.NewUserStore()
	adminUser := user.User{
		Id:       "user1",
		Name:     "Jack",
		Role:     "admin",
		Username: "jacktest",
		Password: "Supersecure",
	}
	basicUser := user.User{
		Id:       "user2",
		Name:     "Test",
		Role:     "basic",
		Username: "stest",
		Password: "Supersecure2",
	}
	userStoreMock.Users = map[string]user.User{
		"user1": adminUser,
		"user2": basicUser,
	}

	jwtStoreMock := jwt_manager.NewJWTManager([]byte("SecretYouShouldHide"))
	testcases := []testcase{
		{
			desc:           "HAPPY admin user is allowed",
			cookie:         cookie(&adminUser, jwtStoreMock),
			expectedStatus: 200,
		},
		{
			desc:           "HAPPY basic user is not allowed",
			cookie:         cookie(&basicUser, jwtStoreMock),
			expectedStatus: 403,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.desc, func(t *testing.T) {
			rMock := mux.NewRouter()

			permissionsMock := NewPermissionsMgr(
				userStoreMock,
				jwtStoreMock,
			)

			rMock.Use(permissionsMock.Permissions)
			rMock.HandleFunc("/test", test()).Methods("GET")

			w := httptest.NewRecorder()

			r, _ := http.NewRequest("GET", "/test", nil)
			r.AddCookie(testCase.cookie)
			rMock.ServeHTTP(w, r)

			assert.Equal(t, w.Result().StatusCode, testCase.expectedStatus)
		})
	}
}

// cookie returns a cookie with JWT based on the user passed and
// an expiration date of now + 10 seconds
func cookie(user *user.User, jwtMgr jwt_manager.JWTManager) *http.Cookie {
	expirationDate := time.Now().Add(10 * time.Second)
	jwt, _ := jwtMgr.CreateJWT(user, expirationDate)

	return &http.Cookie{
		Name:    "token",
		Value:   jwt,
		Expires: expirationDate,
	}
}

// test is a test endpoint set up to test the permissions
func test() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("This is a test"))
	}
}
