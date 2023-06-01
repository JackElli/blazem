package permissions

import (
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/users"

	"net/http"
)

type PermissionsMgr struct {
	UserStore users.UserStorer
	JWTMgr    jwt_manager.JWTManager
}

func NewPermissionsMgr(userStore users.UserStorer, jwtMgr jwt_manager.JWTManager) *PermissionsMgr {
	return &PermissionsMgr{
		UserStore: userStore,
		JWTMgr:    jwtMgr,
	}
}

func (e *PermissionsMgr) Permissions(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("token")
		if err != nil {
			w.WriteHeader(401)
			return
		}

		jwtStr := c.Value
		userId, err := e.JWTMgr.GetCurrentUserId(jwtStr)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		user, err := e.UserStore.Get(userId)
		if err != nil {
			w.WriteHeader(401)
			return
		}
		if user.Role == "admin" {
			h.ServeHTTP(w, req)
			return
		}

		logger.Logger.Info(userId + " tried to do something which needs special permissions and was rejected.")
		w.WriteHeader(403)
	})
}
