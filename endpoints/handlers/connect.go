package handlers

import (
	"blazem/global"
	"blazem/logging"
	"encoding/json"
	"fmt"
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

	ip := req.URL.Query().Get("ip")

	if ip != "" {
		if !global.AlreadyInNodeMap(ip) {
			global.NODE_MAP = append(global.NODE_MAP, &global.Node{Ip: ip, Pinged: time.Now(),
				PingCount: 0, Rank: global.FOLLOWER, Data: sync.Map{}, Active: true,
				RecentQueries: map[string]string{}, Rules: map[string]global.Rule{}})
		} else {
			indexOfNode := global.IndexOfNodeIpInNodeMap(ip)
			global.NODE_MAP[indexOfNode].Active = true
			global.NODE_MAP[indexOfNode].PingCount = 0
		}
		global.Logger.Log(ip+" has connected", logging.GOOD)
	}

	jsonNodeMap, err := json.Marshal(global.NODE_MAP)
	if err != nil {
		fmt.Println("Cannot Marshal json")
		return
	}
	json.NewEncoder(w).Encode(jsonNodeMap)
}
