package query

import (
	types "blazem/domain/endpoint"
	global_types "blazem/domain/global"
	"blazem/global"
	"blazem/query"
	"encoding/json"
	"net/http"
	"time"
)

func NewQueryHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return QueryHandler
}

func QueryHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &QueryEndpoint{
		Endpoint: *e,
	}
	return pe.queryHandler
}

// We fetch the query entered by the user, we send that to JAQL, then
// we send back the results to the client. We also want to add these to
// recent queries so the user can easily get back to queries they've
// previously entered.
func (e *QueryEndpoint) queryHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"query"})
	if req.Method != "GET" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}

	var queryVal = req.URL.Query().Get("query")
	if queryVal == "" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "No query param sent",
		})
		return
	}
	var dataToSend = make([]global_types.SendData, 0)
	query.LoadIntoMemory(*e.Endpoint.Node)

	queryResult, timeTaken, _, errors := query.Execute(queryVal, "")
	if len(errors) != 0 {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Errors found in query response",
			Data: errors,
		})
		return
	}

	for _, res := range queryResult {
		if res["type"] != "text" {
			res["value"] = "file"
		}
		dataJSON, err := json.Marshal(res)
		if err != nil {
			e.Endpoint.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Cannot marshal query data {" + err.Error() + "}",
			})
			return
		}
		var getJSON global.JsonData
		err = json.Unmarshal(dataJSON, &getJSON)
		if err != nil {
			e.Endpoint.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Cannot unmarshal query data {" + err.Error() + "}",
			})
			return
		}
		dataToSend = append(dataToSend, global_types.SendData{
			Key:  getJSON["key"].(string),
			Data: getJSON,
		})
	}
	e.Endpoint.Node.RecentQueries[queryVal] = time.Now().Format("2006-01-02 15:04:05")
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Completed query successfully",
		Data: global_types.SendQueryData{
			Docs:      dataToSend,
			TimeTaken: timeTaken,
		},
	})
}
