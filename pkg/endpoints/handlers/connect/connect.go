package connect

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/logging"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// We need to connect a node to the cluster; we check for ip, if it is already
// in the node map, we set to active (because it must be active as it's sent a
// connect request). If it's not in the nodemap, we add it.
func Connect(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ip := mux.Vars(req)["ip"]
		err := updateNodeMap(ip)
		if err != nil {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully connected",
			Data: global.NODE_MAP,
		})
	}
}

// We need to append a node to the nodemap or, if the node is already in
// the nodemap, we can activate it again
func updateNodeMap(ip string) error {
	if ip == "" {
		return errors.New("IP nothing")
	}
	if !global.AlreadyInNodeMap(ip) {
		global.NODE_MAP = append(global.NODE_MAP, &global.Node{Ip: ip, Pinged: time.Now(),
			PingCount: 0, Rank: global.FOLLOWER, Data: sync.Map{}, Active: true,
			RecentQueries: map[string]string{}, Rules: map[string]global.Rule{}})
		global.Logger.Log(ip+" has connected", logging.GOOD)
		return nil
	}
	indexOfNode := global.IndexOfNodeIpInNodeMap(ip)
	global.NODE_MAP[indexOfNode].Active = true
	global.NODE_MAP[indexOfNode].PingCount = 0
	return nil
}
