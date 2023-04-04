package handlers

import (
	"blazem/global"
	"blazem/query"
	"encoding/json"
	"net/http"
	"time"
)

func QueryHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.queryHandler
}

func (node *Node) queryHandler(w http.ResponseWriter, req *http.Request) {
	// We fetch the query entered by the user, we send that to JAQL, then
	// we send back the results to the client. We also want to add these to
	// recent queries so the user can easily get back to queries they've
	// previously entered.

	WriteHeaders(w, []string{"query"})

	queryVal := req.URL.Query().Get("query")
	if queryVal == "" {
		queryVal = req.Header.Get("query")
	}

	query.LoadIntoMemory(global.Node(*node))

	queryResult, timeTaken, _, _ := query.Execute(queryVal, "")
	dataToSend := make([]SendData, 0)

	for _, res := range queryResult {
		if res["type"] != "text" {
			res["value"] = "file"
		}

		dataJSON, _ := json.Marshal(res)
		var getJSON global.JsonData
		json.Unmarshal(dataJSON, &getJSON)
		dataToSend = append(dataToSend, SendData{getJSON["key"].(string), getJSON})
	}

	node.RecentQueries[queryVal] = time.Now().Format("2006-01-02 15:04:05")
	json.NewEncoder(w).Encode(SendQueryData{dataToSend, timeTaken})
}
