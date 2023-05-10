package stats

import (
	"blazem/pkg/domain/endpoint"
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

// We want to fetch stats based on the OS (as the method changes)
func Stats(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var os = runtime.GOOS
		var stats = osStats[os]()
		json.NewEncoder(w).Encode(stats)
	}
}
