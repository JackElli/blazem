package removenode

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"net/http"
)

func NewRemoveNodeHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return RemoveNodeHandler
}

func RemoveNodeHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &RemoveNodeEndpoint{
		Endpoint: *e,
	}
	return pe.removeNodeHandler
}

// We want to remove a node from the node map (master only). We get the index in
// the node map of the node, then we update the node map (removing the 'to remove' node)
// then we save the changes.
func (e *RemoveNodeEndpoint) removeNodeHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"ip"})
	if req.Method != "DELETE" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	if e.Endpoint.Node.Rank == global.FOLLOWER {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Must be a master node to complete this action",
		})
		return
	}
	nodeIpToRemove := req.URL.Query().Get("ip")
	if nodeIpToRemove == "" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "No IP passed",
		})
		return
	}
	indexOfNode := global.IndexOfNodeIpInNodeMap(nodeIpToRemove)
	if indexOfNode == -1 {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Invalid node",
		})
		return
	}
	global.NODE_MAP = append(global.NODE_MAP[:indexOfNode], global.NODE_MAP[indexOfNode+1:]...)
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully removed node from nodemap",
	})
}
