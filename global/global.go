package global

import (
	"blazem/logging"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func marshalNodeMap(nodeMap []*Node) []*TempNode {
	var SEND_MAP []*TempNode
	for _, node := range NODE_MAP {
		var nodeData = make(map[string]interface{}, 0)
		node.Data.Range(func(key, value any) bool {
			nodeData[key.(string)] = value
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

// PingRetry retries the ping 3 times and if
// afer 3 pings there's no response,
// node is 'paused'
func (n *Node) PingRetry(sendData *bytes.Buffer) bool {

	if n == nil || n.Active == false {
		return false
	}

	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		_, err := http.Post("http://"+n.Ip+"/ping",
			"application/json", sendData)
		if err == nil {
			return true
		}
		Logger.Log("PINGING AGAIN", logging.INFO)
	}
	return false
}

func (node *Node) PingEachConnection(jsonNodeMap []byte) {
	for _, n := range NODE_MAP {
		go func(loopn *Node) {

			//don't ping to itself
			if loopn.Ip == node.Ip {
				return
			}

			//ping non active node every 8 seconds
			if loopn.Active == false {

				//2n and 2n+1 (all cases)
				//maybe a better way would be to send a ping to master
				//if node comes back
				if time.Now().Second()%7 != 0 &&
					time.Now().Second()%8 != 0 {
					return
				}
			}

			sendData := bytes.NewBuffer(jsonNodeMap)

			//ping connection
			Logger.Log("PINGING "+loopn.Ip, logging.INFO)
			_, err := http.Post("http://"+loopn.Ip+"/ping",
				"application/json", sendData)

			//retry logic
			if err != nil {
				if !loopn.PingRetry(sendData) {
					Logger.Log("Cannot connect to "+loopn.Ip,
						logging.WARNING)
					loopn.Active = false
					loopn.PingCount = 0
					return
				}
			}
			//send all data to new joiner
			//TODO check if data is the same on nodes
			if loopn.PingCount == 0 {
				Logger.Log("SENDING MAP TO FIRST JOINER",
					logging.INFO)
				//marshall so we're able to send over TCP
				jsonNodeMap, _ := json.Marshal(marshalNodeMap(NODE_MAP))
				sendData := bytes.NewBuffer(jsonNodeMap)
				_, err = http.Post("http://"+loopn.Ip+"/ping",
					"application/json", sendData)
			}
			//increase connection ping count
			loopn.PingCount++
			Logger.Log("PING RECEIVED FROM "+loopn.Ip,
				logging.INFO)

			if loopn.Active == false {
				loopn.Active = true
			}
		}(n)
	}
	node.PingCount++
}

func (node *Node) Ping() {
	//while true
	for true {
		time.Sleep(4 * time.Second)
		//break out if follower (shouldnt be pinging if follower)
		if node.Rank == FOLLOWER {
			return
		}

		if node.Rank == MASTER {
			Logger.Log(string(node.Rank)+" at "+node.Ip+
				" nodemap: "+strings.Join(GetNodeIps(), " "),
				logging.INFO)
		}

		if len(NODE_MAP) == 1 {
			continue
		}

		//check if the data has changed from the data on the map
		jsonNodeMap := checkIfDataChanged()
		node.PingEachConnection(jsonNodeMap)
	}
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
	Logger.Log("Slow response first check at "+
		fmt.Sprintf("%f", timeSinceLastPingAbs)+"s", logging.WARNING)
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
	node.setToMaster()
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

func (node *Node) setToMaster() {
	node.Rank = MASTER
	// pass all data over
	node.Data = NODE_MAP[0].Data
	node.RecentQueries = NODE_MAP[0].RecentQueries
	node.Rules = NODE_MAP[0].Rules

	waitingTimeStr := strconv.Itoa(int(time.Now().Sub(node.Pinged).Seconds()))
	Logger.Log("IM THE MASTER NOW, COPIED ALL DATA FROM PREVIOUS MASTER!!! after waiting for "+waitingTimeStr+"s", logging.GOOD)
	// update node map
	NODE_MAP = NODE_MAP[1:]
	NODE_MAP[0] = node
	// save a backup
	// node.SaveBackup()
	// start pinging again
	go node.Ping()
}

type TempNode struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          Rank
	Data          map[string]interface{}
	Active        bool
	RecentQueries map[string]string //time
	Rules         map[string]Rule
}

// this needs improving, need to check data not just endpoint
func checkIfDataChanged() []byte {
	var jsonNodeMap []byte
	if DataChanged {

		// need to change from syncMap
		// to normal map to send
		jsonNodeMap, _ = json.Marshal(marshalNodeMap(NODE_MAP))
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
