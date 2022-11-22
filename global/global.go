package global

import (
	"bytes"
	"distributed_servers/logging"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Rank string

const (
	MASTER   Rank = "MASTER"
	FOLLOWER Rank = "FOLLOWER"
)

// global vars will clean up
var PORT_START = 3100
var NODE_MAP []*Node
var Logger logging.Logger
var DataChanged bool = false

type Node struct {
	Ip        string
	Pinged    time.Time
	PingCount int
	Rank      Rank
	Data      map[string]string
	Active    bool
}

func (n *Node) PingRetry() bool {
	Logger.Log("PINGING AGAIN", logging.INFO)
	_, err := net.DialTimeout("tcp", n.Ip, 2000*time.Millisecond)
	if err != nil {
		Logger.Log("Cannot connect to "+n.Ip, logging.WARNING)
		n.Active = false
		n.PingCount = 0
		return false
	}
	return true
}

func (node *Node) PingEachConnection(jsonNodeMap []byte) {
	for _, n := range NODE_MAP {
		go func(loopn *Node) {

			//ping non active node every 8 seconds
			if loopn.Active == false {
				//2n and 2n+1 (all cases)
				if time.Now().Second()%7 != 0 && time.Now().Second()%8 != 0 {
					return
				}
			}

			//don't ping to itself
			if loopn.Ip == node.Ip {
				return
			}
			//send all data to new joiner
			if loopn.PingCount == 0 {
				Logger.Log("SENDING MAP TO FIRST JOINER", logging.INFO)
				//marshall so we're able to send over TCP
				jsonNodeMap, _ = json.Marshal(NODE_MAP)
			}
			sendData := bytes.NewBuffer(jsonNodeMap)

			//ping connection
			Logger.Log("PINGING "+loopn.Ip, logging.INFO)
			_, err := net.DialTimeout("tcp", loopn.Ip, 1500*time.Millisecond)

			//increase connection ping count
			loopn.PingCount++

			//retry logic
			if err != nil {
				if !loopn.PingRetry() {
					return
				}
			}
			_, err = http.Post("http://"+loopn.Ip+"/ping", "application/json", sendData)
			Logger.Log("PING RECEIVED FROM "+loopn.Ip, logging.INFO)
			if loopn.Active == false {
				loopn.Active = true
				// //-1 just in case
				// loopn.PingCount = -1
			}
		}(n)
	}
	node.PingCount++
}

func (node *Node) Ping() {
	//while true
	for true {

		//break out if follower (shouldnt be pinging if follower)
		if node.Rank == FOLLOWER {
			return
		}

		time.Sleep(4 * time.Second)

		if node.Rank == MASTER {
			Logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+strings.Join(GetNodeIps(), " "), logging.INFO)
		}

		if len(NODE_MAP) == 1 {
			continue
		}

		if node.Rank == FOLLOWER {
			continue
		}

		//check if the data has changed from the data on the map
		jsonNodeMap := checkIfDataChanged()
		node.PingEachConnection(jsonNodeMap)
	}
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

func (node *Node) CheckForNoPingFromMaster() {

	//master shouldnt be checking
	if node.Rank == MASTER {
		return
	}
	//wait for 5 seconds after ping
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs := time.Now().Sub(node.Pinged).Seconds()
	// if there has been a ping in that time, good!
	if timeSinceLastPingAbs < 1 {
		return
	}
	Logger.Log("Slow response first check at "+fmt.Sprintf("%f", timeSinceLastPingAbs)+"s", logging.WARNING)
	//if not, check for retry ping
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs = time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 8.2 {
		return
	}
	//if no pings in that time, master is down
	Logger.Log("NO PING FROM MASTER!!!", logging.INFO)

	//if node is not next in line, break out
	if !node.isNextInLine() {
		return
	}
	//set that node to master
	node.Rank = MASTER
	node.Data = NODE_MAP[0].Data
	waitingTimeStr := strconv.Itoa(int(time.Now().Sub(node.Pinged).Seconds()))
	Logger.Log("IM THE MASTER NOW, COPIED ALL DATA FROM PREVIOUS MASTER!!! after waiting for "+waitingTimeStr+"s", logging.GOOD)
	//update node map
	NODE_MAP = NODE_MAP[1:]
	NODE_MAP[0] = node
	//start pinging again
	go node.Ping()
}

// this needs improving, need to check data not just endpoint
func checkIfDataChanged() []byte {
	var jsonNodeMap []byte
	if DataChanged {
		Logger.Log("DATA CHANGED", logging.INFO)
		jsonNodeMap, _ = json.Marshal(NODE_MAP)
		DataChanged = false
	} else {
		jsonNodeMap, _ = json.Marshal(getNodeMapWithoutData())
	}
	return jsonNodeMap
}

func getNodeMapWithoutData() []*Node {
	var newmap []*Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank, map[string]string{}, n.Active})
	}
	return newmap
}

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

func GetAllDataToPrint(data map[string]string) []string {
	var retdata []string
	for v := range data {
		retdata = append(retdata, data[v])
	}
	return retdata
}
