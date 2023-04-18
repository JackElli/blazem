package handlers

import (
	"blazem/global"
	"net/http"
)

// Return the results of the nodemap to the client
func NodeMapHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, []string{"all"})

	if req.Method != "GET" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	var nodeMap = make([]WebNodeMap, 0)
	for _, n := range global.NODE_MAP {
		nodeMap = append(nodeMap, WebNodeMap{n.Ip, n.Active})
	}
	JsonResponse(w, EndpointResponse{
		200,
		"Successfully retrieved nodemap",
		nodeMap,
	})
}
