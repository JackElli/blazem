package nodemap

import (
	types "blazem/domain/endpoint"
	"blazem/domain/global"
	"net/http"
)

func NewNodeMapHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return NodeMapHandler
}

func NodeMapHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	ne := &NodeMapEndpoint{
		Endpoint: *e,
	}
	return ne.NodeMapHandler
}

// Return the results of the nodemap to the client
func (e *NodeMapEndpoint) NodeMapHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"all"})
	if req.Method != "GET" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	nodeMap := make([]types.WebNodeMap, 0)
	for _, n := range global.NODE_MAP {
		nodeMap = append(nodeMap, types.WebNodeMap{
			Ip:     n.Ip,
			Active: n.Active,
		})
	}
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully retrieved nodemap",
		Data: nodeMap,
	})
}
