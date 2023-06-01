package recentquery

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"net/http"

	"github.com/gorilla/mux"
)

var route = "/recentQueries"

type RecentQueryMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
}

func NewRecentQueryMgr(router *mux.Router, node *node.Node, responder responder.Responder) *RecentQueryMgr {
	return &RecentQueryMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
	}
}

// Returns a list of recently entered queries
func (e *RecentQueryMgr) RecentQuery() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		dataToSend := e.Node.RecentQueries
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg:  "Successfully retrieved recent queries",
			Data: dataToSend,
		})
	}
}

func (e *RecentQueryMgr) Register() {
	e.Router.HandleFunc(route, e.RecentQuery()).Methods("GET")
}
