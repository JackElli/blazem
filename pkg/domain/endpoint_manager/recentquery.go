package endpoint_manager

import (
	types "blazem/pkg/domain/endpoint"
	"net/http"
)

// Returns a list of recently entered queries
func (e *EndpointManager) RecentQuery() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		dataToSend := e.Node.RecentQueries
		// if len(dataToSend) == 0 {
		// 	e.Responder.Respond(w, types.EndpointResponse{
		// 		Code: 200,
		// 		Msg:  "Successfully retrieved recent queries",
		// 		Data: []uint8{},
		// 	})
		// 	return
		// }
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg:  "Successfully retrieved recent queries",
			Data: dataToSend,
		})
	}
}
