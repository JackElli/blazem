package query

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/query"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// We fetch the query entered by the user, we send that to JAQL, then
// we send back the results to the client. We also want to add these to
// recent queries so the user can easily get back to queries they've
// previously entered.
func Query(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		queryVal := mux.Vars(req)["query"]
		if queryVal == "" {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No query param sent",
			})
			return
		}
		dataToSend := make([]types.SendData, 0)
		query.LoadIntoMemory(*r.Node)
		queryResult, timeTaken, _, errors := query.Execute(queryVal, "")

		if len(errors) != 0 {
			r.Respond(w, types.EndpointResponse{
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
				r.Respond(w, types.EndpointResponse{
					Code: 500,
					Msg:  "Cannot marshal query data {" + err.Error() + "}",
				})
				return
			}
			var getJSON global.JsonData
			err = json.Unmarshal(dataJSON, &getJSON)
			if err != nil {
				r.Respond(w, types.EndpointResponse{
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
		r.Node.RecentQueries[queryVal] = time.Now().Format("2006-01-02 15:04:05")
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Completed query successfully",
			Data: types.SendQueryData{
				Docs:      dataToSend,
				TimeTaken: timeTaken,
			},
		})
	}
}
