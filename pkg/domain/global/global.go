package global

import (
	"blazem/pkg/logging"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PingRetry retries the ping 3 times and if afer 3 pings there's no response,
// node is 'paused'
func (n *Node) PingRetry(sendData *bytes.Buffer) bool {
	if n == nil || n.Active == false {
		return false
	}
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		_, err := http.Post("http://"+n.Ip+"/ping", "application/json", sendData)
		if err == nil {
			return true
		}
		Logger.Log("PINGING AGAIN", logging.INFO)
	}
	return false
}

// We want to ping each follower node to make sure they know, the
// master is still active
func (node *Node) PingEachConnection(jsonNodeMap []byte) {
	for _, n := range NODE_MAP {
		go func(loopn *Node) {
			if loopn.Ip == node.Ip {
				return
			}
			if loopn.Active == false {
				if time.Now().Second()%7 != 0 && time.Now().Second()%8 != 0 {
					return
				}
			}
			sendData := bytes.NewBuffer(jsonNodeMap)
			Logger.Log("PINGING "+loopn.Ip, logging.INFO)

			_, err := http.Post("http://"+loopn.Ip+"/ping", "application/json", sendData)
			if err != nil {
				if !loopn.PingRetry(sendData) {
					Logger.Log("Cannot connect to "+loopn.Ip,
						logging.WARNING)
					loopn.Active = false
					loopn.PingCount = 0
					return
				}
			}
			if loopn.PingCount == 0 {
				Logger.Log("SENDING MAP TO FIRST JOINER", logging.INFO)
				jsonNodeMap, _ := json.Marshal(MarshalNodeMap(NODE_MAP))
				sendData := bytes.NewBuffer(jsonNodeMap)
				_, err = http.Post("http://"+loopn.Ip+"/ping", "application/json", sendData)
			}
			loopn.PingCount++
			Logger.Log("PING RECEIVED FROM "+loopn.Ip, logging.INFO)
			if loopn.Active == false {
				loopn.Active = true
			}
		}(n)
	}
	node.PingCount++
}

// Every for seconds, we want to ping each connection
func (node *Node) Ping() {
	for true {
		time.Sleep(4 * time.Second)
		if node.Rank == FOLLOWER {
			return
		}
		if len(NODE_MAP) == 1 {
			continue
		}
		Logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+strings.Join(GetNodeIps(), " "),
			logging.INFO)
		jsonNodeMap := checkIfDataChanged()
		node.PingEachConnection(jsonNodeMap)
	}
}

// We want to check if the master is still alive
func (node *Node) CheckForNoPingFromMaster() {
	if node.Rank == MASTER {
		return
	}
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs := time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 1 {
		return
	}
	Logger.Log("Slow response first check at "+fmt.Sprintf("%f", timeSinceLastPingAbs)+"s",
		logging.WARNING)
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs = time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 8.2 {
		return
	}
	Logger.Log("NO PING FROM MASTER!!!", logging.INFO)
	if node.isNextInLine() {
		node.setToMaster()
	}
}

// Set this node to master status and put all 'replicas' to 'active'
func (node *Node) setToMaster() {
	node.Rank = MASTER
	node.Data = NODE_MAP[0].Data
	node.RecentQueries = NODE_MAP[0].RecentQueries
	node.Rules = NODE_MAP[0].Rules

	waitingTimeStr := strconv.Itoa(int(time.Now().Sub(node.Pinged).Seconds()))
	Logger.Log("IM THE MASTER NOW, COPIED ALL DATA FROM PREVIOUS MASTER!!! after waiting for "+waitingTimeStr+"s", logging.GOOD)
	NODE_MAP = NODE_MAP[1:]
	NODE_MAP[0] = node
	go node.Ping()
}

// reads from data storage puts all docs to memory on load
func (node *Node) ReadFromLocal() {
	files, _ := ioutil.ReadDir("data/")
	if len(files) == 0 {
		return
	}
	for _, file := range files {
		key := file.Name()
		data, _ := ioutil.ReadFile("data/" + key)
		var dataJSON JsonData
		json.Unmarshal(data, &dataJSON)
		node.Data.Store(key, (Document)(dataJSON))
	}
	Logger.Log("Loaded files into memory.", logging.INFO)
}

// We want to send data across nodes
func MarshalNodeMap(nodeMap []*Node) []*TempNode {
	var SEND_MAP []*TempNode
	for _, node := range NODE_MAP {
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
			node.Ip,
			node.Pinged,
			node.PingCount,
			node.Rank,
			nodeData,
			node.Active,
			node.RecentQueries,
			node.Rules,
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

// We want to write a document to disk
func WriteDocToDisk(value Document) {
	dataToWrite, _ := json.Marshal(value)
	path := "data/"
	_ = os.MkdirAll(path, os.ModePerm)
	os.WriteFile("data/"+value["key"].(string), []byte(dataToWrite), os.ModePerm)
}

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
