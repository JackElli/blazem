package global

import (
	"encoding/json"
	"strconv"
	"sync"
)

// return the ips stored in the nodemap
func GetNodeIps() []string {
	var nodeips []string

	for _, n := range NODE_MAP {
		node := n.Ip + ":" + strconv.FormatBool(n.Active)
		nodeips = append(nodeips, node)
	}
	return nodeips
}

func AlreadyInNodeMap(ip string) bool {
	for _, n := range NODE_MAP {
		if n.Ip == ip {
			return true
		}
	}
	return false
}

func IndexOfNodeIpInNodeMap(ip string) int {
	for i, n := range NODE_MAP {
		if n.Ip == ip {
			return i
		}
	}
	return -1
}

func (node *Node) isNextInLine() bool {
	//get next true value
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

// this needs improving, need to check data not just endpoint
func checkIfDataChanged() []byte {
	var jsonNodeMap []byte
	if DataChanged {
		// need to change from syncMap
		// to normal map to send
		jsonNodeMap, _ = json.Marshal(MarshalNodeMap(NODE_MAP))
		DataChanged = false
	} else {
		jsonNodeMap, _ = json.Marshal(getNodeMapWithoutData())
	}
	return jsonNodeMap
}

func getNodeMapWithoutData() []*Node {
	var newmap []*Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank,
			sync.Map{}, n.Active, n.RecentQueries, n.Rules})
	}
	return newmap
}
