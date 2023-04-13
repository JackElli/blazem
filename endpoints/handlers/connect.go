package handlers

import (
	"blazem/global"
	"blazem/logging"
	"errors"
	"net/http"
	"sync"
	"time"
)

func ConnectHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.connectHandler
}

func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {
	// We need to connect a node to the cluster; we check for ip, if it is already
	// in the node map, we set to active (because it must be active as it's sent a
	// connect request). If it's not in the nodemap, we add it.
	WriteHeaders(w, []string{"ip"})

	if req.Method != "POST" {
		JsonResponse(w, EndpointResponse{
			500,
			"Should not be getting",
			nil,
		})
		return
	}

	var ip = req.URL.Query().Get("ip")
	var err = updateNodeMap(ip)

	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			err.Error(),
			nil,
		})
		return
	}

	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			"Cannot marshal nodemap",
			nil,
		})
		return
	}
	JsonResponse(w, EndpointResponse{
		200,
		"Successfully connected",
		global.NODE_MAP,
	})
}

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
