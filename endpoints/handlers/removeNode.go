package handlers

import (
	"blazem/global"
	"blazem/logging"
	"encoding/json"
	"net/http"
)

func RemoveNodeHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.removeNodeHandler
}

func (node *Node) removeNodeHandler(w http.ResponseWriter, req *http.Request) {
	// we want to remove a node from the node map (master only). We get the index in
	// the node map of the node, then we update the node map (removing the 'to remove' node)
	// then we save the changes.
	WriteHeaders(w, []string{"ip"})

	if node.Rank == global.FOLLOWER {
		return
	}

	nodeIpToRemove := req.URL.Query().Get("ip")
	if nodeIpToRemove == "" {
		nodeIpToRemove = req.Header.Get("ip")
	}

	indexOfNode := global.IndexOfNodeIpInNodeMap(nodeIpToRemove)
	if indexOfNode == -1 {
		return
	}
	global.Logger.Log("removed node: "+nodeIpToRemove+" from the nodemap", logging.GOOD)

	global.NODE_MAP = append(global.NODE_MAP[:indexOfNode], global.NODE_MAP[indexOfNode+1:]...)
	json.NewEncoder(w).Encode("removed node")

}
