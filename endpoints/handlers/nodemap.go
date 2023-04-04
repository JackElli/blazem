package handlers

import (
	"blazem/global"
	"encoding/json"
	"net/http"
)

func NodeMapHandler(w http.ResponseWriter, req *http.Request) {
	// Return the results of the nodemap to the client
	WriteHeaders(w, []string{"all"})
	nodeMapResp := []WebNodeMap{}
	for _, n := range global.NODE_MAP {
		nodeMapResp = append(nodeMapResp, WebNodeMap{n.Ip, n.Active})
	}
	json.NewEncoder(w).Encode(nodeMapResp)
}
