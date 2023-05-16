package recentquery

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"net/http"
)

// Returns a list of recently entered queries
func RecentQuery(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		dataToSend := e.Node.RecentQueries
		if len(dataToSend) == 0 {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 200,
				Msg:  "Successfully retrieved recent queries",
				Data: []uint8{},
			})
			return
		}
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved recent queries",
			Data: dataToSend,
		})
	}
}
