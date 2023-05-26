package endpoint_manager

import (
	"blazem/pkg/domain/endpoint"
	"net/http"
)

func (e *EndpointManager) NodeMap() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		nodeMap := make([]endpoint.WebNodeMap, 0)
		for _, n := range e.Node.NodeMap {
			nodeMap = append(nodeMap, endpoint.WebNodeMap{
				Ip:     n.Ip,
				Active: n.Active,
				Rank:   n.Rank,
			})
		}

		e.Responder.Respond(w, 200, endpoint.EndpointResponse{
			Msg:  "Successfully retrieved nodemap",
			Data: nodeMap,
		})
	}
}
