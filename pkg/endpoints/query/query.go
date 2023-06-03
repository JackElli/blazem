package query

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/query"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var route = "/query"

type QueryMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
	_Query    *query.Query
}

func NewQueryMgr(router *mux.Router, node *node.Node, responder responder.Responder, query *query.Query) *QueryMgr {
	return &QueryMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
		_Query:    query,
	}
}

// We fetch the query entered by the user, we send that to JAQL, then
// we send back the results to the client. We also want to add these to
// recent queries so the user can easily get back to queries they've
// previously entered.
func (e *QueryMgr) Query() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var queryVal struct {
			Query string `json:"query"`
		}
		json.NewDecoder(req.Body).Decode(&queryVal)

		if queryVal.Query == "" {
			e.Responder.Error(w, 500, errors.New("No query param sent"))
			return
		}

		dataToSend := make([]types.SendData, 0)
		queryResult, timeTaken, _, errs := e._Query.Execute(queryVal.Query)

		if len(errs) != 0 {
			e.Responder.Error(w, 500, errors.New("Errors found in query response"))
			return
		}

		for _, res := range queryResult {
			if res["type"] != "text" {
				res["value"] = "file"
			}
			dataJSON, err := json.Marshal(res)
			if err != nil {
				e.Responder.Error(w, 500, err)
				return
			}
			var getJSON global.JsonData
			err = json.Unmarshal(dataJSON, &getJSON)
			if err != nil {
				e.Responder.Error(w, 500, err)
				return
			}
			dataToSend = append(dataToSend, types.SendData{
				Key:  getJSON["key"].(string),
				Data: getJSON,
			})
		}
		e.Node.RecentQueries[queryVal.Query] = time.Now().Format("2006-01-02 15:04:05")
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg: "Completed query successfully",
			Data: types.SendQueryData{
				Docs:      dataToSend,
				TimeTaken: timeTaken,
			},
		})
	}
}

func (e *QueryMgr) Register() {
	e.Router.HandleFunc(route, e.Query()).Methods("POST")
}
