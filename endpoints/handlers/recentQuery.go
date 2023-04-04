package handlers

import (
	"encoding/json"
	"net/http"
)

func GetRecentQueriesHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.getRecentQueriesHandler
}

func (node *Node) getRecentQueriesHandler(w http.ResponseWriter, req *http.Request) {
	// Returns a list of recently entered queries
	WriteHeaders(w, []string{})

	dataToSend := node.RecentQueries
	if len(dataToSend) == 0 {
		json.NewEncoder(w).Encode([]uint8{})
		return
	}
	json.NewEncoder(w).Encode(dataToSend)
}
