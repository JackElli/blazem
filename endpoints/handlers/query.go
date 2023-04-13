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

	if req.Method != "GET" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	var queryVal = req.URL.Query().Get("query")
	if queryVal == "" {
		JsonResponse(w, EndpointResponse{
			500,
			"No query param sent",
			nil,
		})
		return
	}
	var dataToSend = make([]SendData, 0)
	query.LoadIntoMemory(global.Node(*node))

	queryResult, timeTaken, _, errors := query.Execute(queryVal, "")
	if len(errors) != 0 {
		JsonResponse(w, EndpointResponse{
			500,
			"Errors found in query response",
			errors,
		})
		return
	}

	for _, res := range queryResult {
		if res["type"] != "text" {
			res["value"] = "file"
		}
		var dataJSON, err = json.Marshal(res)
		if err != nil {
			JsonResponse(w, EndpointResponse{
				500,
				"Cannot marshal query data {" + err.Error() + "}",
				nil,
			})
			return
		}

		var getJSON global.JsonData
		err = json.Unmarshal(dataJSON, &getJSON)
		if err != nil {
			JsonResponse(w, EndpointResponse{
				500,
				"Cannot unmarshal query data {" + err.Error() + "}",
				nil,
			})
			return
		}
		dataToSend = append(dataToSend, SendData{getJSON["key"].(string), getJSON})
	}
	node.RecentQueries[queryVal] = time.Now().Format("2006-01-02 15:04:05")
	JsonResponse(w, EndpointResponse{
		200,
		"Completed query successfully",
		SendQueryData{dataToSend, timeTaken},
	})
}
