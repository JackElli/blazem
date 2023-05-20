package query

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"encoding/json"
	"net/http"
	"time"
)

// We fetch the query entered by the user, we send that to JAQL, then
// we send back the results to the client. We also want to add these to
// recent queries so the user can easily get back to queries they've
// previously entered.
func Query(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var queryVal struct {
			Query string `json:"query"`
		}
		json.NewDecoder(req.Body).Decode(&queryVal)

		if queryVal.Query == "" {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No query param sent",
			})
			return
		}

		dataToSend := make([]types.SendData, 0)
		queryResult, timeTaken, _, errors := e.Query.Execute(queryVal.Query)

		if len(errors) != 0 {
			e.Responder.Respond(w, types.EndpointResponse{
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
				e.Responder.Respond(w, types.EndpointResponse{
					Code: 500,
					Msg:  "Cannot marshal query data {" + err.Error() + "}",
				})
				return
			}
			var getJSON global.JsonData
			err = json.Unmarshal(dataJSON, &getJSON)
			if err != nil {
				e.Responder.Respond(w, types.EndpointResponse{
					Code: 500,
					Msg:  "Cannot unmarshal query data {" + err.Error() + "}",
				})
				return
			}
			dataToSend = append(dataToSend, types.SendData{
				Key:  getJSON["key"].(string),
				Data: getJSON,
			})
		}
		e.Node.RecentQueries[queryVal.Query] = time.Now().Format("2006-01-02 15:04:05")
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Completed query successfully",
			Data: types.SendQueryData{
				Docs:      dataToSend,
				TimeTaken: timeTaken,
			},
		})
	}
}
