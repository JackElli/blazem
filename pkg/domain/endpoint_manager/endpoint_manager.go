package endpoint_manager

import (
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
)

type EndpointManager struct {
	Node      *node.Node
	Responder *responder.Respond
}

func NewEndpointManager(node *node.Node, responder *responder.Respond) *EndpointManager {
	return &EndpointManager{
		Node:      node,
		Responder: responder,
	}
}
