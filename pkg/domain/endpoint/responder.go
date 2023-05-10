package endpoint

import (
	"blazem/pkg/domain/global"
	"encoding/json"
	"net/http"
	"strings"
)

type Respond struct {
	Node *global.Node
}

type Responder interface {
	Respond(w http.ResponseWriter, response EndpointResponse)
	WriteHeaders(w http.ResponseWriter, extras []string)
}

func NewResponder(node *global.Node) *Respond {
	return &Respond{
		Node: node,
	}
}

// Send a JSON response back
func (r *Respond) Respond(w http.ResponseWriter, response EndpointResponse) {
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}

// We want to write headers for each request, the content type and
// the CORS settings
func (r *Respond) WriteHeaders(w http.ResponseWriter, extras []string) {
	extra := strings.Join(extras, ",")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, "+extra)
	w.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET")
}
