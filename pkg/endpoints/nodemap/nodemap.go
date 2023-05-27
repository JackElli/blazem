package nodemap

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/users"
	"net/http"

	"github.com/gorilla/mux"
)

type NodemapMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
	UserStore users.UserStorer
}

func NewNodemapMgr(router *mux.Router, node *node.Node, responder responder.Responder) *NodemapMgr {
	return &NodemapMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
	}
}

func (e *NodemapMgr) NodeMap() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		nodeMap := make([]endpoint.WebNodeMap, 0)
		for _, n := range e.Node.NodeMap {
			nodeMap = append(nodeMap, endpoint.WebNodeMap{
				Ip:     n.Ip,
				Active: n.Active,
				Rank:   n.Rank,
			})
		}

		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg:  "Successfully retrieved nodemap",
			Data: nodeMap,
		})
	}
}

func (e *NodemapMgr) Register() {
	e.Router.HandleFunc("/nodemap", e.NodeMap()).Methods("GET")
}
