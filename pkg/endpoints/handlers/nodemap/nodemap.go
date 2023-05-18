package nodemap

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"net/http"
)

func NodeMap(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		nodeMap := make([]endpoint.WebNodeMap, 0)
		for _, n := range e.Node.NodeMap {
			nodeMap = append(nodeMap, endpoint.WebNodeMap{
				Ip:     n.Ip,
				Active: n.Active,
				Rank:   n.Rank,
			})
		}
		e.Responder.Respond(w, endpoint.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved nodemap",
			Data: nodeMap,
		})
	}
}
