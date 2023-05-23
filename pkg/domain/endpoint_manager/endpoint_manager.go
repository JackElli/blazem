package endpoint_manager

import (
	"blazem/pkg/domain/jwt_manager"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/storer"
	"blazem/pkg/domain/users"
	"blazem/pkg/query"
)

type EndpointManager struct {
	Node       *node.Node
	Responder  responder.Responder
	JWTManager jwt_manager.JWTManager
	Query      *query.Query
	UserStore  users.UserStorer
	DataStore  storer.Storer
}

func NewEndpointManager(node *node.Node, responder responder.Responder, jwtMgr jwt_manager.JWTManager, query *query.Query, userStore users.UserStorer, store storer.Storer) *EndpointManager {
	return &EndpointManager{
		Node:       node,
		Responder:  responder,
		JWTManager: jwtMgr,
		Query:      query,
		UserStore:  userStore,
		DataStore:  store,
	}
}
