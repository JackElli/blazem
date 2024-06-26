package main

import (
	"blazem/pkg/domain/global"
	blazem_node "blazem/pkg/domain/node"
	"sync"
	"time"
)

// We want to initialise this node.
//
// Create a new node with default settings, set up the logger and pick the correct IP and port.
//
// We want to set up the endpoints so that the UI can call it and append this node to the NODE_MAP to be passed around each node.
//
// If we are the master, we set master attributes. We also want to load from disk to fetch any data saved locally.
func main() {
	node := &blazem_node.Node{
		Ip:            "",
		Pinged:        time.Now(),
		PingCount:     0,
		Rank:          global.FOLLOWER,
		Data:          sync.Map{},
		Active:        true,
		RecentQueries: map[string]string{},
		Rules:         map[string]global.Rule{},
	}

	RunSetup(node)

	for true {
		time.Sleep(150 * time.Millisecond)
	}
}
