package endpoint

import (
	"blazem/pkg/domain/global"
	"encoding/json"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}
