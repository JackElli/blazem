package endpoint

import (
	"blazem/domain/global"
	"encoding/json"
	"net/http"
	"strings"
)

type EndpointType string

const (
	ASYNC EndpointType = "async"
	SYNC  EndpointType = "sync"
)

type Endpoint struct {
	Node        *global.Node
	Route       string
	Handler     func(*Endpoint) func(w http.ResponseWriter, req *http.Request)
	Description string
	Type        EndpointType
}

type EndpointResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type IEndpoint interface {
	Respond(w http.ResponseWriter, response EndpointResponse)
	WriteHeaders(w http.ResponseWriter, extras []string)
}

// Send a JSON response back
func (e *Endpoint) Respond(w http.ResponseWriter, response EndpointResponse) {
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}

// We want to write headers for each request, the content type and
// the CORS settings
func (e *Endpoint) WriteHeaders(w http.ResponseWriter, extras []string) {
	extra := strings.Join(extras, ",")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, "+extra)
}