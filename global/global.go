package global

import (
	"blazem/logging"
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

func MarshalNodeMap(nodeMap []*Node) []*TempNode {
	// We want to send data across nodes
	var SEND_MAP []*TempNode
	for _, node := range NODE_MAP {
		var nodeData = make(map[string]interface{}, 0)
		files, _ := ioutil.ReadDir("data/")

		node.Data.Range(func(key, value any) bool {
			var docKey = key.(string)
			var jsonData map[string]interface{}
			if value.(map[string]interface{})["type"] != "text" {
				if len(files) == 0 {
					return true
				}
				for _, file := range files {
					var key = file.Name()
					var data, _ = ioutil.ReadFile("data/" + key)
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

func UnmarshalNodeMap(nodeMap []*TempNode) []*Node {
	// The opposite of Marshal, for retrieving data from nodes
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

func (n *Node) PingRetry(sendData *bytes.Buffer) bool {
	// PingRetry retries the ping 3 times and if afer 3 pings there's no response,
	// node is 'paused'

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

func (node *Node) PingEachConnection(jsonNodeMap []byte) {
	// We want to ping each follower node to make sure they know, the master is still active
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
			var sendData = bytes.NewBuffer(jsonNodeMap)

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

func (node *Node) Ping() {
	// Every for seconds, we want to ping each connection
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

		var jsonNodeMap = checkIfDataChanged()
		node.PingEachConnection(jsonNodeMap)
	}
}

func (node *Node) CheckForNoPingFromMaster() {
	// We want to check if the master is still alive
	if node.Rank == MASTER {
		return
	}

	time.Sleep(4100 * time.Millisecond)
	var timeSinceLastPingAbs = time.Now().Sub(node.Pinged).Seconds()
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

func (node *Node) setToMaster() {
	// Set this node to master status and put all 'replicas' to 'active'
	node.Rank = MASTER
	node.Data = NODE_MAP[0].Data
	node.RecentQueries = NODE_MAP[0].RecentQueries
	node.Rules = NODE_MAP[0].Rules

	var waitingTimeStr = strconv.Itoa(int(time.Now().Sub(node.Pinged).Seconds()))
	Logger.Log("IM THE MASTER NOW, COPIED ALL DATA FROM PREVIOUS MASTER!!! after waiting for "+waitingTimeStr+"s", logging.GOOD)

	NODE_MAP = NODE_MAP[1:]
	NODE_MAP[0] = node
	go node.Ping()
}

func (node *Node) ReadFromLocal() {
	// reads from data storage puts all docs to memory on load
	var files, _ = ioutil.ReadDir("data/")
	if len(files) == 0 {
		return
	}
	for _, file := range files {
		var key = file.Name()
		var data, _ = ioutil.ReadFile("data/" + key)
		var dataJSON JsonData
		json.Unmarshal(data, &dataJSON)
		node.Data.Store(key, (Document)(dataJSON))
	}
	Logger.Log("Loaded files into memory.", logging.INFO)
}

func WriteDocToDisk(value Document) {
	// We want to write a document to disk
	dataToWrite, _ := json.Marshal(value)
	path := "data/"
	_ = os.MkdirAll(path, os.ModePerm)
	os.WriteFile("data/"+value["key"].(string), []byte(dataToWrite), os.ModePerm)
}
