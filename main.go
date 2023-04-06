package main

import (
	"blazem/endpoints"
	"blazem/global"
	"blazem/query"
	"sync"
	"time"
)

type Node global.Node

func main() {
	// We want to initialise this node. Create a new node with default settings
	// set up the logger and pick the correct IP and port. We want to set up the endpoints
	// so that the UI can call it and append this node to the global NODE_MAP to be passed
	// around each node. If we are the master, we set master attributes. We also
	// want to load from disk to fetch any data saved locally. Then, for JAQL we need to load
	// our query JSON into memory.
	var node Node = Node{
		Ip:            "",
		Pinged:        time.Now(),
		PingCount:     0,
		Rank:          global.FOLLOWER,
		Data:          sync.Map{},
		Active:        true,
		RecentQueries: map[string]string{},
		Rules:         map[string]global.Rule{},
	}

	var masterip string = ""
	var localip = getLocalIp()

	global.GlobalNode = (*global.Node)(&node)
	setupLogger()

	go node.pickPort(localip)
	endpoints.SetupEndpoints((*global.Node)(&node))
	global.NODE_MAP = append(global.NODE_MAP, (*global.Node)(&node))

	if masterip == node.Ip {
		node.setNodeMasterAttrs()
	}

	(*global.Node)(&node).ReadFromLocal()
	go (*global.Node)(&node).Ping()

	// // run the rule checker
	// go (*endpoints.Node)(&node).CheckRules()

	query.LoadIntoMemory(global.Node(node))

	for true {
		time.Sleep(150 * time.Millisecond)
	}
}
