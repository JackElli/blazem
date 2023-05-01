package stats

import (
	types "blazem/pkg/domain/endpoint"
	"encoding/json"
	"net/http"
	"runtime"
)

var osStats = map[string]func() types.Stats{
	"windows": nil,
	"linux":   nil,
	"darwin":  nil,
}

func NewStatsHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return StatsHandler
}

func StatsHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &StatsEndpoint{
		Endpoint: *e,
	}
	return pe.statsHandler
}

// We want to fetch stats based on the OS (as the method changes)
func (e *StatsEndpoint) statsHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, nil)
	var os = runtime.GOOS
	var stats = osStats[os]()
	json.NewEncoder(w).Encode(stats)
}
