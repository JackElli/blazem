package nodes

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"net/http"

	"github.com/gorilla/mux"
)

// We want to remove a node from the node map (master only). We get the index in
// the node map of the node, then we update the node map (removing the 'to remove' node)
// then we save the changes.
func RemoveNode(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank == global.FOLLOWER {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Must be a master node to complete this action",
			})
			return
		}
		nodeIpToRemove := mux.Vars(req)["ip"]
		if nodeIpToRemove == "" {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No IP passed",
			})
			return
		}
		indexOfNode := global.IndexOfNodeIpInNodeMap(nodeIpToRemove)
		if indexOfNode == -1 {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Invalid node",
			})
			return
		}
		global.NODE_MAP = append(global.NODE_MAP[:indexOfNode], global.NODE_MAP[indexOfNode+1:]...)
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully removed node from nodemap",
		})
	}
}
