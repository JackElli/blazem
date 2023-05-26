package endpoint_manager

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// We want to remove a node from the node map (master only). We get the index in
// the node map of the node, then we update the node map (removing the 'to remove' node)
// then we save the changes.
func (e *EndpointManager) RemoveNode() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank == global.FOLLOWER {
			e.Responder.Error(w, 500, errors.New("Must be a master node to complete this action"))
			return
		}
		nodeIpToRemove := mux.Vars(req)["ip"]
		indexOfNode := e.Node.IndexOfNodeIpInNodeMap(nodeIpToRemove)
		if indexOfNode == -1 {
			e.Responder.Error(w, 500, errors.New("Invalid node"))
			return
		}

		e.Node.NodeMap = append(e.Node.NodeMap[:indexOfNode], e.Node.NodeMap[indexOfNode+1:]...)
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg: "Successfully removed node from nodemap",
		})
	}
}
