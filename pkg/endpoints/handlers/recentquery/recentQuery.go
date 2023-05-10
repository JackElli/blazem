package recentquery

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"net/http"
)

// Returns a list of recently entered queries
func RecentQuery(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
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
