package responder

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"encoding/json"
	"net/http"
)

type Respond struct {
	Node *global.Node
}

type Responder interface {
	Respond(w http.ResponseWriter, response endpoint.EndpointResponse)
}

// NewResponder returns a respond struct implementing the
// responder interface
func NewResponder(node *global.Node) *Respond {
	return &Respond{
		Node: node,
	}
}

// Send a JSON response back
func (r *Respond) Respond(w http.ResponseWriter, response endpoint.EndpointResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}
