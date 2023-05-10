package nodemap

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"net/http"
)

func NodeMap(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		nodeMap := make([]types.WebNodeMap, 0)
		for _, n := range global.NODE_MAP {
			nodeMap = append(nodeMap, types.WebNodeMap{
				Ip:     n.Ip,
				Active: n.Active,
			})
		}
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successfully retrieved nodemap",
			Data: nodeMap,
		})
	}
}
