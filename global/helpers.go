package global

import (
	"encoding/json"
	"strconv"
	"sync"
)

func GetNodeIps() []string {
	// Return the ips stored in the nodemap
	var nodeips []string
	for _, n := range NODE_MAP {
		node := n.Ip + ":" + strconv.FormatBool(n.Active)
		nodeips = append(nodeips, node)
	}
	return nodeips
}

func AlreadyInNodeMap(ip string) bool {
	// Return true if this node is already in the cluster
	for _, n := range NODE_MAP {
		if n.Ip == ip {
			return true
		}
	}
	return false
}

func IndexOfNodeIpInNodeMap(ip string) int {
	// Get the index of the node in the cluster
	for i, n := range NODE_MAP {
		if n.Ip == ip {
			return i
		}
	}
	return -1
}

func (node *Node) isNextInLine() bool {
	// Get next true value
	for _, n := range NODE_MAP {
		if n.Active == false {
			continue
		}
		if n.Ip == node.Ip {
			return true
		}
	}
	return false
}

func checkIfDataChanged() []byte {
	// This checks and returns the data on the node if it has changed
	var jsonNodeMap []byte
	if DataChanged {
		jsonNodeMap, _ = json.Marshal(MarshalNodeMap(NODE_MAP))
		DataChanged = false
	} else {
		jsonNodeMap, _ = json.Marshal(getNodeMapWithoutData())
	}
	return jsonNodeMap
}

func getNodeMapWithoutData() []*Node {
	// Returns the data without any data for easy sending
	var newmap []*Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank,
			sync.Map{}, n.Active, n.RecentQueries, n.Rules})
	}
	return newmap
}
