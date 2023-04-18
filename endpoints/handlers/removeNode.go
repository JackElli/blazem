package handlers

import (
	"blazem/global"
	"net/http"
)

func RemoveNodeHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.removeNodeHandler
}

// We want to remove a node from the node map (master only). We get the index in
// the node map of the node, then we update the node map (removing the 'to remove' node)
// then we save the changes.
func (node *Node) removeNodeHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, []string{"ip"})

	if req.Method != "DELETE" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	if node.Rank == global.FOLLOWER {
		JsonResponse(w, EndpointResponse{
			500,
			"Must be a master node to complete this action",
			nil,
		})
		return
	}

	var nodeIpToRemove = req.URL.Query().Get("ip")
	if nodeIpToRemove == "" {
		JsonResponse(w, EndpointResponse{
			500,
			"No IP passed",
			nil,
		})
		return
	}

	var indexOfNode = global.IndexOfNodeIpInNodeMap(nodeIpToRemove)
	if indexOfNode == -1 {
		JsonResponse(w, EndpointResponse{
			500,
			"Invalid node",
			nil,
		})
		return
	}
	global.NODE_MAP = append(global.NODE_MAP[:indexOfNode], global.NODE_MAP[indexOfNode+1:]...)
	JsonResponse(w, EndpointResponse{
		200,
		"Successfully removed node from nodemap",
		nil,
	})
}
