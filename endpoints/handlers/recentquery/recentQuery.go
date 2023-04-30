package recentquery

import (
	types "blazem/domain/endpoint"
	"net/http"
)

func NewRecentQueryHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return RecentQueryHandler
}

func RecentQueryHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &RecentQueryEndpoint{
		Endpoint: *e,
	}
	return pe.getRecentQueriesHandler
}

// Returns a list of recently entered queries
func (e *RecentQueryEndpoint) getRecentQueriesHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{})

	if req.Method != "GET" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}

	dataToSend := e.Endpoint.Node.RecentQueries
	if len(dataToSend) == 0 {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved recent queries",
			Data: []uint8{},
		})
		return
	}
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully retrieved recent queries",
		Data: dataToSend,
	})
}
