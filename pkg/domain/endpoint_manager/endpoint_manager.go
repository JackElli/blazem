package endpoint_manager

import (
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"blazem/pkg/query"
)

type EndpointManager struct {
	Node      *node.Node
	Query     *query.Query
	Responder responder.Responder
	Store     storer.Storer
}

func NewEndpointManager(node *node.Node, responder responder.Responder, query *query.Query, store storer.Storer) *EndpointManager {
	return &EndpointManager{
		Node:      node,
		Responder: responder,
		Query:     query,
		Store:     store,
	}
}
