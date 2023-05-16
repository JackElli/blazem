package recentquery

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/responder"
	"net/http"
)

// Returns a list of recently entered queries
func RecentQuery(r *responder.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		dataToSend := r.Node.RecentQueries
		if len(dataToSend) == 0 {
			r.Respond(w, types.EndpointResponse{
				Code: 200,
				Msg:  "Successfully retrieved recent queries",
				Data: []uint8{},
			})
			return
		}
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved recent queries",
			Data: dataToSend,
		})
	}
}
