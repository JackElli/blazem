package handlers

import (
	"net/http"
)

func GetRecentQueriesHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.getRecentQueriesHandler
}

func (node *Node) getRecentQueriesHandler(w http.ResponseWriter, req *http.Request) {
	// Returns a list of recently entered queries
	WriteHeaders(w, []string{})

	if req.Method != "GET" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	var dataToSend = node.RecentQueries
	if len(dataToSend) == 0 {
		JsonResponse(w, EndpointResponse{
			200,
			"Successfully retrieved recent queries",
			[]uint8{},
		})
		return
	}
	JsonResponse(w, EndpointResponse{
		200,
		"Successfully retrieved recent queries",
		dataToSend,
	})
}
