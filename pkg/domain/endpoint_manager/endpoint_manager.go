package endpoint_manager

import (
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"blazem/pkg/domain/users"
	"blazem/pkg/query"
)

type EndpointManager struct {
	Node      *node.Node
	Query     *query.Query
	Responder responder.Responder
	UserStore users.UserStorer
	DataStore storer.Storer
}

func NewEndpointManager(node *node.Node, responder responder.Responder, query *query.Query, userStore users.UserStorer, store storer.Storer) *EndpointManager {
	return &EndpointManager{
		Node:      node,
		Responder: responder,
		Query:     query,
		UserStore: userStore,
		DataStore: store,
	}
}
