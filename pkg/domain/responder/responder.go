package responder

import (
	"blazem/pkg/domain/endpoint"
	"encoding/json"
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter, response endpoint.EndpointResponse)
}

type Respond struct{}

// NewResponder returns a respond struct implementing the
// responder interface
func NewResponder() *Respond {
	return &Respond{}
}

// Send a JSON response back
func (r *Respond) Respond(w http.ResponseWriter, response endpoint.EndpointResponse) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(response.Code)

	json.NewEncoder(w).Encode(response)
}
