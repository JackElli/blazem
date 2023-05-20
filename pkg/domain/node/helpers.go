package node

import (
	"blazem/pkg/domain/global"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync"
)

// Return the ips stored in the nodemap
func (node *Node) GetNodeIps() []string {
	var nodeips []string
	for _, n := range node.NodeMap {
		node := n.Ip + ":" + strconv.FormatBool(n.Active)
		nodeips = append(nodeips, node)
	}
	return nodeips
}

// Return true if this node is already in the cluster
func (node *Node) AlreadyInNodeMap(ip string) bool {
	for _, n := range node.NodeMap {
		if n.Ip == ip {
			return true
		}
	}
	return false
}

// Get the index of the node in the cluster
func (node *Node) IndexOfNodeIpInNodeMap(ip string) int {
	for i, n := range node.NodeMap {
		if n.Ip == ip {
			return i
		}
	}
	return -1
}

// Returns the data without any data for easy sending
func (node *Node) GetNodeMapWithoutData() []*Node {
	var newmap []*Node
	for _, n := range node.NodeMap {
		newmap = append(newmap, &Node{
			n.Ip,
			n.Pinged,
			0,
			n.Rank,
			sync.Map{},
			n.Active,
			n.RecentQueries,
			n.Rules,
			n.NodeMap,
			n.UserStore,
		})
	}
	return newmap
}

// We want to send data across nodes
func (node *Node) MarshalNodeMap() []*TempNode {
	var SEND_MAP []*TempNode
	for _, node := range node.NodeMap {
		nodeData := make(map[string]interface{}, 0)
		files, _ := ioutil.ReadDir("data/")
		node.Data.Range(func(key, value any) bool {
			docKey := key.(string)
			var jsonData map[string]interface{}
			if value.(map[string]interface{})["type"] != "text" {
				if len(files) == 0 {
					return true
				}
				for _, file := range files {
					key := file.Name()
					data, _ := ioutil.ReadFile("data/" + key)
					if key != docKey {
						continue
					}
					json.Unmarshal(data, &jsonData)
					nodeData[docKey] = jsonData
					return true
				}
			}
			nodeData[docKey] = value
			return true
		})
		tempNode := TempNode{
			Ip:            node.Ip,
			Pinged:        node.Pinged,
			PingCount:     node.PingCount,
			Rank:          node.Rank,
			Data:          nodeData,
			Active:        node.Active,
			RecentQueries: node.RecentQueries,
			Rules:         node.Rules,
		}
		SEND_MAP = append(SEND_MAP, &tempNode)
	}
	return SEND_MAP
}

// The opposite of Marshal, for retrieving data from nodes
func UnmarshalNodeMap(nodeMap []*TempNode) []*Node {
	var SEND_MAP []*Node
	for _, node := range nodeMap {
		var nodeData sync.Map
		for key, value := range node.Data {
			nodeData.Store(key, value)
		}
		SEND_MAP = append(SEND_MAP, &Node{
			Ip:            node.Ip,
			Pinged:        node.Pinged,
			PingCount:     node.PingCount,
			Rank:          node.Rank,
			Data:          nodeData,
			Active:        node.Active,
			RecentQueries: node.RecentQueries,
			Rules:         node.Rules,
		})
	}
	return SEND_MAP
}

// Get next true value
func (node *Node) IsNextInLine() bool {
	for _, n := range node.NodeMap {
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
func (node *Node) CheckIfDataChanged() []byte {
	var jsonNodeMap []byte
	if global.DataChanged {
		jsonNodeMap, _ = json.Marshal(node.MarshalNodeMap())
		global.DataChanged = false
	} else {
		jsonNodeMap, _ = json.Marshal(node.GetNodeMapWithoutData())
	}
	return jsonNodeMap
}
