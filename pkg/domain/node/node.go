package node

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/users"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type INode interface {
	CheckForNoPingFromMaster()
	Ping()
	PingEachConnection(jsonNodeMap []byte)
	PingRetry(sendData *bytes.Buffer) bool
	ReadFromLocal()
}

type Node struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          global.Rank
	Data          sync.Map
	Active        bool
	RecentQueries map[string]string //time
	Rules         map[string]global.Rule
	NodeMap       []*Node
	UserStore     *users.UserStore
}
type TempNode struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          global.Rank
	Data          map[string]interface{}
	Active        bool
	RecentQueries map[string]string //time
	Rules         map[string]global.Rule
}

type SetupStep struct {
	Description string
	Fn          func() error
}

var GlobalNode *Node

func NewNode() *Node {
	return &Node{
		Ip:            "",
		Pinged:        time.Now(),
		PingCount:     0,
		Rank:          global.FOLLOWER,
		Data:          sync.Map{},
		Active:        true,
		RecentQueries: map[string]string{},
		Rules:         map[string]global.Rule{},
	}
}

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
		logger.Logger.Info("Pinging again.")
	}
	return false
}

// We want to ping each follower node to make sure they know, the
// master is still active
func (node *Node) PingEachConnection(jsonNodeMap []byte) {
	for _, n := range node.NodeMap {
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
			logger.Logger.Info("Pinging " + loopn.Ip)

			_, err := http.Post("http://"+loopn.Ip+"/ping", "application/json", sendData)
			if err != nil {
				if !loopn.PingRetry(sendData) {
					logger.Logger.Warn("Cannot connect to " + loopn.Ip)
					loopn.Active = false
					loopn.PingCount = 0
					return
				}
			}
			if loopn.PingCount == 0 {
				logger.Logger.Info("Sending map to first joiner.")
				jsonNodeMap, _ := json.Marshal(node.MarshalNodeMap())
				sendData := bytes.NewBuffer(jsonNodeMap)
				_, err = http.Post("http://"+loopn.Ip+"/ping", "application/json", sendData)
			}
			loopn.PingCount++
			logger.Logger.Info("Ping received from " + loopn.Ip)
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
		if node.Rank == global.FOLLOWER {
			return
		}
		if len(node.NodeMap) == 1 {
			continue
		}
		jsonNodeMap := node.CheckIfDataChanged()
		node.PingEachConnection(jsonNodeMap)
	}
}

// We want to check if the master is still alive
func (node *Node) CheckForNoPingFromMaster() {
	if node.Rank == global.MASTER {
		return
	}
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs := time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 1 {
		return
	}
	logger.Logger.Warn("Slow response first check at " + fmt.Sprintf("%f", timeSinceLastPingAbs) + "s")
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs = time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 8.2 {
		return
	}
	logger.Logger.Info("No ping from master.")
	if node.IsNextInLine() {
		node.SetToMaster()
	}
}

// Set this node to master status and put all 'replicas' to 'active'
func (node *Node) SetToMaster() {
	node.Rank = global.MASTER
	node.Data = node.NodeMap[0].Data
	node.RecentQueries = node.NodeMap[0].RecentQueries
	node.Rules = node.NodeMap[0].Rules

	waitingTimeStr := strconv.Itoa(int(time.Now().Sub(node.Pinged).Seconds()))
	logger.Logger.Info("I'm the master now. I've copied all of the data from the previous master after waiting for " + waitingTimeStr + "s")
	node.NodeMap = node.NodeMap[1:]
	node.NodeMap[0] = node
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
		var dataJSON global.JsonData
		json.Unmarshal(data, &dataJSON)
		node.Data.Store(key, (global.Document)(dataJSON))
	}
	logger.Logger.Info("Loaded files into memory.")
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

// We want to write a document to disk
func (node *Node) WriteDocToDisk(value global.Document) {
	dataToWrite, _ := json.Marshal(value)
	path := "data/"
	_ = os.MkdirAll(path, os.ModePerm)
	os.WriteFile("data/"+value["key"].(string), []byte(dataToWrite), os.ModePerm)
}

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
