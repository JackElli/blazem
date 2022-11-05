package main

import (
	"bytes"
	"distributed_servers/logging"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var PORT_START = 3100
var NODE_MAP []*Node

var dataChanged bool = false

type Node struct {
	Ip        string
	Pinged    time.Time
	PingCount int
	Rank      Rank
	Data      map[string]string
}

// utils
// return the index of the node in the nodemap
func indexOfNodeInNodeMap(node *Node) int {
	for i, n := range NODE_MAP {
		if n.Ip == node.Ip {
			return i
		}
	}
	return -1
}

// return the data stored in the nodemap
func getNodeDatas() []map[string]string {
	var nodedata []map[string]string
	for _, n := range NODE_MAP {
		nodedata = append(nodedata, n.Data)
	}
	return nodedata
}

// return the data stored in the nodemap
func getNodeIps() []string {
	var nodeips []string
	for _, n := range NODE_MAP {
		nodeips = append(nodeips, n.Ip)
	}
	return nodeips
}

// check if data in transit has changed from the node data
func hasDataChanged(d1 map[string]string, d2 map[string]string) bool {
	if len(d1) != len(d2) {
		// fmt.Print("LENGTH DIFF")
		return true
	}
	for _, k := range d1 {
		if v, exists := d2[k]; exists {
			if v != d1[k] {
				return true
			}
			return false
		}
	}
	return false
}

func getNodeMapWithoutData() []*Node {
	var newmap []*Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank, map[string]string{}})
	}
	return newmap
}

func nodeMapPointertoMem() []Node {
	var newmap []Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, Node{n.Ip, n.Pinged, n.PingCount, n.Rank, n.Data})
	}
	return newmap
}
func getAllDataToPrint(data map[string]string) string {
	retdata := ""
	for v := range data {
		retdata += data[v] + ","
	}
	return retdata
}
func getLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// ping connections
func (node *Node) tryListen(ip string) {

	//wait for 500 milliseconds
	time.Sleep(500 * time.Millisecond)

	//listen on selected port
	portstr := ip
	if strings.Count(ip, ":") > 1 {
		portstr = strings.Split(ip, ":")[0]
	}

	logging.Log("trying on "+portstr, logging.INFO)
	l, err := net.Listen("tcp", portstr)

	//if theres an error in connecting, stop
	if err != nil {
		return
	}

	//if there's no error set the nodes port to the current port
	node.Ip = ip
	logging.Log("connected on "+portstr, logging.GOOD)

	//serve http requests on this port
	http.Serve(l, nil)
}

func (node *Node) pickPort(ip string) {

	//for each node, try to listen
	for i := 0; i < 3; i++ {
		node.tryListen(ip + ":" + strconv.Itoa(PORT_START+i))

		//if theres an error, break out
		if node.Ip != "" {
			break
		}
	}
}

func (node *Node) ping() {

	//while true
	for true {

		//print the rank of the node and wait for 2 secs
		logging.Log(string(node.Rank)+" at "+node.Ip, logging.INFO)
		// ,
		time.Sleep(2 * time.Second)

		//only send ping if its a master
		if node.Rank == MASTER {
			// fmt.Println("Data stored: ", node.Data)
			jsonNodeMap, _ := json.Marshal(NODE_MAP)
			//check if the data has changed from the data on the map
			if dataChanged {
				fmt.Println("DATA CHANGED")
				dataChanged = false
			} else {
				jsonNodeMap, _ = json.Marshal(getNodeMapWithoutData())
			}

			for _, n := range NODE_MAP {

				//don't ping to itself
				if n.Ip != node.Ip {
					//marshall so we're able to send over TCP
					if n.PingCount == 0 {
						logging.Log("SENDING MAP TO FIRST JOINER", logging.INFO)
						jsonNodeMap, _ = json.Marshal(NODE_MAP)
					}
					sendData := bytes.NewBuffer(jsonNodeMap)
					//set timeout to 2 seconds
					c := &http.Client{
						Timeout: 2 * time.Second,
					}
					logging.Log("PINGING "+"http://"+n.Ip, logging.INFO)
					resp, err := c.Post("http://"+n.Ip+"/ping", "application/json", sendData)
					n.PingCount++

					if err != nil {
						logging.Log(err.Error(), logging.WARNING)
						indexOfNode := indexOfNodeInNodeMap(n)
						NODE_MAP = append(NODE_MAP[:indexOfNode], NODE_MAP[indexOfNode+1:]...)
						continue
					}
					if resp.Header.Get("pinged") == "true" {
						logging.Log("PING RECEIVED FROM "+n.Ip, logging.GOOD)
					}
				}
			}

			// if node.PingCount%5 == 0 {
			// 	rn := rand.Intn(100)
			// 	newdata := make(map[string]string)
			// 	for i := 0; i < rn; i++ {
			// 		r := byte(rand.Intn(40))
			// 		newdata[strconv.Itoa(i)] = string(r)
			// 	}
			// 	node.Data = newdata
			// 	// fmt.Println("NEWDATA ", nodeD)
			// }
			node.PingCount++

		}
	}
}

func (node *Node) checkForPingFromMaster() {

	//while true
	for true {

		//if its master, break
		if node.Rank == MASTER {
			return
		}

		//if the ping hasnt changed (eg ping didnt set it to true)
		// fmt.Println(time.Now().Sub(node.Pinged).Seconds())
		if time.Now().Sub(node.Pinged).Seconds() > 2.5 {

			//node must be down
			logging.Log("NO PING FROM MASTER, MUST BE DOWN!!!", logging.INFO)

			//if node is next in line to throne
			//NEED TO CHECK IF NODE MAP IS CORRECT
			if node.Ip == NODE_MAP[1].Ip {

				//set that node to master
				node.Rank = MASTER
				node.Data = NODE_MAP[0].Data

				logging.Log("IM THE MASTER NOW, COPIED ALL DATA FROM PREVIOUS MASTER!!!", logging.GOOD)

			} else {
				node.Pinged = time.Now()
				//sub 2 seconds from this
			}

			//update node map
			NODE_MAP = NODE_MAP[1:]
			NODE_MAP[0] = node
			// fmt.Println("UPDATING NODE MAP", getNodeIps())
		}
	}
}

// handlers
func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {

	//only receive ping if its a follower
	if node.Rank == FOLLOWER {

		//fetch data of ping (nodemap)
		body, _ := ioutil.ReadAll(req.Body)

		//move nodemap to local memory
		var localnm []*Node
		json.Unmarshal(body, &localnm)

		//add the changed node map
		if len(NODE_MAP) != len(localnm) {
			NODE_MAP = append(NODE_MAP, localnm[len(localnm)-1])
		}

		if len(localnm[0].Data) == 0 {

			// fmt.Println("data is empty, must not have changed")
			node.Pinged = time.Now()
			// fmt.Println("PINGED")
			w.Header().Add("pinged", "true")
			return
		}

		logging.Log("UPDATED DATA ON THIS NODE!", logging.GOOD)

		NODE_MAP = []*Node{}
		for _, j := range localnm {
			NODE_MAP = append(NODE_MAP, j)
		}

		//it has been successfully pinged
		node.Pinged = time.Now()
		w.Header().Add("pinged", "true")
	}

}
func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {

	//find the port of client trying to connect
	ip := req.URL.Query().Get("ip")

	//if port isnt 0 (there is a valid port)
	if ip != "" {
		//add to the node map
		logging.Log(ip+" has connected", logging.GOOD)
		NODE_MAP = append(NODE_MAP, &Node{ip, time.Now(), 0, FOLLOWER, map[string]string{}})
	}

	//write back to the client with the rank of the node its trying to connect to
	w.Header().Add("rank", string(node.Rank))
	// fmt.Println("SEDNING node map", NODE_MAP)
	jsonNodeMap, _ := json.Marshal(NODE_MAP)
	sendData := bytes.NewBuffer(jsonNodeMap)
	w.Write(sendData.Bytes())
}

func (node *Node) dataHandler(w http.ResponseWriter, req *http.Request) {
	if node.Rank == MASTER {

		if req.Method == "POST" {

			//TODO send multiple key and values

			var dataToSet []string
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &dataToSet)

			setKey := dataToSet[0]
			setVal := dataToSet[1]
			node.Data[setKey] = setVal

			dataChanged = true
			w.Header().Set("response", "done")
			return
		}

		datagetkey := req.Header.Get("key")

		if datagetkey != "" {
			dataToSend := node.Data[datagetkey]
			w.Header().Set("val", dataToSend)
		}

	}
}

// main func
func main() {

	//init node and set to follower (true until proven otherwise)
	var node Node
	node.Rank = FOLLOWER

	//set ips
	localip := getLocalIp()
	connectoip := ""
	masterip := ""

	fmt.Println("Type in ip to connect, or enter if master.")
	fmt.Scanf("%s", &connectoip)
	if connectoip == "" {
		//get ip of machine
		//THIS IS THE IP TO CONNECT TO
		connectoip = localip
		masterip = localip + ":3100"
	} else {

		//connect to server
		resp, err := http.Get("http://" + connectoip + "/connect")
		if err == nil && resp.Header.Get("rank") == "MASTER" {
			masterip = connectoip
		} else {
			logging.Log("THAT IS NOT A MASTER", logging.ERROR)
			return
		}
	}

	//pick a port for this node to be on you local ip
	go node.pickPort(localip)
	logging.Log("The master is "+masterip, logging.INFO)

	//handle incoming connection
	http.HandleFunc("/connect", node.connectHandler)
	http.HandleFunc("/ping", node.pingHandler)
	http.HandleFunc("/data", node.dataHandler)

	//this needs to be changed, basically, wait for port selection to be done
	time.Sleep(3 * time.Second)
	//add to local nodemap (will be replicated if its master)
	NODE_MAP = append(NODE_MAP, &node)

	//if this node is not on the master port, then its a follower
	if masterip == node.Ip {
		node.Rank = MASTER
		node.Data = map[string]string{}
		node.Data["testkey"] = "hello this is a test"
		node.Data["testkey2"] = "hello this is a second test"

	}

	//connect to port
	if node.Rank == FOLLOWER {
		logging.Log("TRYING TO CONNECT TO "+"http://"+masterip, logging.INFO)
		_, err := http.Get("http://" + *&masterip + "/connect?ip=" + node.Ip)
		if err != nil {
			logging.Log(err.Error(), logging.WARNING)
		}
	}

	//ping handling
	node.Pinged = time.Now()
	go node.ping()
	go node.checkForPingFromMaster()

	for true {
	}
}

//TODO

// NEED TO WAIT FOR MASTER TO COME ONLINE
// GET DATA FROM MASTER (for some api or sdk)
// SOMETIMES 2 MASTERS are ONLINE AT ONCE
