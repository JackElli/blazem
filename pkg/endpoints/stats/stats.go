package stats

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
)

type StatsMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
}

func NewStatsMgr(router *mux.Router, node *node.Node, responder responder.Responder) *StatsMgr {
	return &StatsMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
	}
}

var osStats = map[string]func() types.Stats{
	"windows": nil,
	"linux":   nil,
	"darwin":  nil,
}

// We want to fetch stats based on the OS (as the method changes)
func (e *StatsMgr) Stats() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var os = runtime.GOOS
		var stats = osStats[os]()
		json.NewEncoder(w).Encode(stats)
	}
}

func (e *StatsMgr) Register() {
	e.Router.HandleFunc("/stats", e.Stats()).Methods("GET")
}
