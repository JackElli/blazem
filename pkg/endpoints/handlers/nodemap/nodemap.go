package nodemap

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"net/http"
)

func NodeMap(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		nodeMap := make([]types.WebNodeMap, 0)
		for _, n := range global.NODE_MAP {
			nodeMap = append(nodeMap, types.WebNodeMap{
				Ip:     n.Ip,
				Active: n.Active,
			})
		}
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved nodemap",
			Data: nodeMap,
		})
	}
}
