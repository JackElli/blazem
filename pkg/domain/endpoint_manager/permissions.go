package endpoint_manager

import (
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/middleware"
	"net/http"
)

func (em *EndpointManager) Permissions(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("token")
		if err != nil {
			w.WriteHeader(401)
			return
		}

		jwtStr := c.Value
		userId, err := middleware.GetCurrentUserId(jwtStr)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		user, err := em.Node.UserStore.Get(userId)
		if err != nil {
			w.WriteHeader(401)
			return
		}
		if user.Role == "admin" {
			logger.Logger.Info(userId + " tried to do something which needs special permissions and they are admin.")
			h.ServeHTTP(w, req)
			return
		}

		logger.Logger.Info(userId + " tried to do something which needs special permissions and was rejected.")
		w.WriteHeader(403)
	})
}
