package endpoint_manager

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/responder"
	"blazem/pkg/domain/users"
)

type EndpointManager struct {
	Node      *global.Node
	Responder *responder.Respond
	UserStore *users.UserStore
}

func NewEndpointManager(node *global.Node, responder *responder.Respond, userStore *users.UserStore) *EndpointManager {
	return &EndpointManager{
		Node:      node,
		Responder: responder,
		UserStore: userStore,
	}
}
