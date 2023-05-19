package endpoint_manager

import (
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/query"
)

type EndpointManager struct {
	Node      *node.Node
	Query     *query.Query
	Responder *responder.Respond
}

func NewEndpointManager(node *node.Node, responder *responder.Respond, query *query.Query) *EndpointManager {
	return &EndpointManager{
		Node:      node,
		Responder: responder,
		Query:     query,
	}
}
