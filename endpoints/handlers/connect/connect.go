package connect

import (
	types "blazem/domain/endpoint"
	"blazem/global"
	"blazem/logging"
	"errors"
	"net/http"
	"sync"
	"time"
)

func NewConnectHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return ConnectHandler
}

func ConnectHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	ce := &ConnectEndpoint{
		Endpoint: *e,
	}
	return ce.connectHandler
}

// We need to connect a node to the cluster; we check for ip, if it is already
// in the node map, we set to active (because it must be active as it's sent a
// connect request). If it's not in the nodemap, we add it.
func (e *ConnectEndpoint) connectHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"ip"})
	if req.Method != "POST" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	ip := req.URL.Query().Get("ip")
	err := updateNodeMap(ip)
	if err != nil {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully connected",
		Data: global.NODE_MAP,
	})
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
