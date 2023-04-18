package global

import (
	"encoding/json"
	"strconv"
	"sync"
)

// Return the ips stored in the nodemap
func GetNodeIps() []string {
	var nodeips []string
	for _, n := range NODE_MAP {
		node := n.Ip + ":" + strconv.FormatBool(n.Active)
		nodeips = append(nodeips, node)
	}
	return nodeips
}

// Return true if this node is already in the cluster
func AlreadyInNodeMap(ip string) bool {
	for _, n := range NODE_MAP {
		if n.Ip == ip {
			return true
		}
	}
	return false
}

// Get the index of the node in the cluster
func IndexOfNodeIpInNodeMap(ip string) int {
	for i, n := range NODE_MAP {
		if n.Ip == ip {
			return i
		}
	}
	return -1
}

// Get next true value
func (node *Node) isNextInLine() bool {
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

// This checks and returns the data on the node if it has changed
func checkIfDataChanged() []byte {
	var jsonNodeMap []byte
	if DataChanged {
		jsonNodeMap, _ = json.Marshal(MarshalNodeMap(NODE_MAP))
		DataChanged = false
	} else {
		jsonNodeMap, _ = json.Marshal(getNodeMapWithoutData())
	}
	return jsonNodeMap
}

// Returns the data without any data for easy sending
func getNodeMapWithoutData() []*Node {
	var newmap []*Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank,
			sync.Map{}, n.Active, n.RecentQueries, n.Rules})
	}
	return newmap
}
