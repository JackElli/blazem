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
	_, err := net.DialTimeout("tcp", n.Ip, 2000*time.Millisecond)
	if err != nil {
		Logger.Log("Cannot connect to "+n.Ip, logging.WARNING)
		n.Active = false
		return false
	}
	return true
}

func (node *Node) PingEachConnection(jsonNodeMap []byte) {
	for _, n := range NODE_MAP {

		//don't ping to itself
		if n.Ip == node.Ip || n.Active == false {
			continue
		}
		//marshall so we're able to send over TCP
		//send all data to new joiner
		if n.PingCount == 0 {
			Logger.Log("SENDING MAP TO FIRST JOINER", logging.INFO)
			jsonNodeMap, _ = json.Marshal(NODE_MAP)
		}
		sendData := bytes.NewBuffer(jsonNodeMap)
		//ping connection
		Logger.Log("PINGING "+n.Ip, logging.INFO)
		_, err := net.DialTimeout("tcp", n.Ip, 1500*time.Millisecond)
		//increase connection ping count
		n.PingCount++
		//retry logic
		if err != nil {
			if !n.PingRetry() {
				continue
			}
		}
		Logger.Log("PING RECEIVED FROM "+n.Ip, logging.INFO)
		_, err = http.Post("http://"+n.Ip+"/ping", "application/json", sendData)
		// p.Close()
	}
	node.PingCount++
}

func (node *Node) Ping() {

	if node.Rank == FOLLOWER {
		return
	}
	//while true
	for true {

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

func (node *Node) CheckForNoPingFromMaster() {

	//THIS NEEDS TO BE BETTER

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
	time.Sleep(4 * time.Second)
	timeSinceLastPingAbs = time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 8 {
		return
	}
	//if no pings in that time, master is down
	Logger.Log("NO PING FROM MASTER!!!", logging.INFO)

	//if node is not next in line, break out
	if node.Ip != NODE_MAP[1].Ip {
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
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank, map[string]string{}, true})
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
