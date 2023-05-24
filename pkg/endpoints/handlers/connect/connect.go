package connect

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/logger"
	blazem_node "blazem/pkg/domain/node"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// We need to connect a node to the cluster; we check for ip, if it is already
// in the node map, we set to active (because it must be active as it's sent a
// connect request). If it's not in the nodemap, we add it.
func Connect(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ip := mux.Vars(req)["ip"]
		err := updateNodeMap(e.Node, ip)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg:  "Successfully connected",
			Data: e.Node.NodeMap,
		})
	}
}

// We need to append a node to the nodemap or, if the node is already in
// the nodemap, we can activate it again
func updateNodeMap(node *blazem_node.Node, ip string) error {
	if ip == "" {
		return errors.New("IP not entered")
	}
	if !node.AlreadyInNodeMap(ip) {
		node.NodeMap = append(node.NodeMap, &blazem_node.Node{
			Ip:            ip,
			Pinged:        time.Now(),
			PingCount:     0,
			Rank:          global.FOLLOWER,
			Data:          sync.Map{},
			Active:        true,
			RecentQueries: map[string]string{},
			Rules:         map[string]global.Rule{},
			NodeMap:       node.NodeMap,
		})
		logger.Logger.Info(ip + " has connected")
		return nil
	}
	indexOfNode := node.IndexOfNodeIpInNodeMap(ip)
	node.NodeMap[indexOfNode].Active = true
	node.NodeMap[indexOfNode].PingCount = 0
	return nil
}
