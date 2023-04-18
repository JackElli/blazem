package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
)

var osStats = map[string]func() Stats{
	"windows": getWindowsStats,
	"linux":   getLinuxStats,
	"darwin":  getMacStats,
}

func StatsHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.statsHandler
}

// We want to fetch stats based on the OS (as the method changes)
func (node *Node) statsHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, nil)
	var os = runtime.GOOS
	var stats = osStats[os]()
	json.NewEncoder(w).Encode(stats)
}
