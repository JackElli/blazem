package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
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

var currentmasterip string = ""

type Node struct {
	Ip        string
	Pinged    time.Time
	PingCount int
	Rank      Rank
	Data      []byte
}

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
func getNodeDatas() [][]byte {
	var nodedata [][]byte
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

// need to connect to master
func (node *Node) tryListen(ip string) {

	//wait for 500 milliseconds
	time.Sleep(500 * time.Millisecond)

	//listen on selected port
	portstr := ip
	if strings.Count(ip, ":") > 1 {
		portstr = strings.Split(ip, ":")[0]
	}

	fmt.Println("trying on " + portstr)
	l, err := net.Listen("tcp", portstr)

	//if theres an error in connecting, stop
	if err != nil {
		return
	}

	//if there's no error set the nodes port to the current port
	node.Ip = ip
	fmt.Println("connected on " + portstr)

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

// check if data in transit has changed from the node data
func hasDataChanged(d1 []byte, d2 []byte) bool {
	if len(d1) != len(d2) {
		return true
	}
	for i, d := range d1 {
		if d != d2[i] {
			return true
		}
	}
	return false
}

func getNodeMapWithoutData() []*Node {
	var newmap []*Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank, []byte{}})
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

func (node *Node) ping() {

	var previoussend []Node = nodeMapPointertoMem()
	//while true
	for true {

		//print the rank of the node and wait for 2 secs
		fmt.Println(node.Rank)
		// ,
		time.Sleep(2 * time.Second)
		//only send ping if its a master
		if node.Rank == MASTER {
			// fmt.Println("Data stored: ", node.Data)
			jsonNodeMap, _ := json.Marshal(NODE_MAP)
			//check if the data has changed from the data on the map
			// fmt.Println(previoussend, node.Data)
			if !hasDataChanged(previoussend[0].Data, node.Data) {
				jsonNodeMap, _ = json.Marshal(getNodeMapWithoutData())
			} else {
				previoussend = nodeMapPointertoMem()
			}

			for _, n := range NODE_MAP {

				//don't ping to itself
				if n.Ip != node.Ip {
					//marshall so we're able to send over TCP
					if n.PingCount == 0 {
						fmt.Println("SENDING MAP TO FIRST JOINER")
						jsonNodeMap, _ = json.Marshal(NODE_MAP)
					}
					sendData := bytes.NewBuffer(jsonNodeMap)
					fmt.Println("PINGING "+"http://"+n.Ip, " NODEMAP ", getNodeIps())
					resp, err := http.Post("http://"+n.Ip+"/ping", "application/json", sendData)
					n.PingCount++

					if err != nil {
						fmt.Println(err)
						indexOfNode := indexOfNodeInNodeMap(n)
						NODE_MAP = append(NODE_MAP[:indexOfNode], NODE_MAP[indexOfNode+1:]...)
						continue
					}
					if resp.Header.Get("pinged") == "true" {
						fmt.Println("PING RECEIVED FROM ", n.Ip)
					}
				}
			}

			if node.PingCount%10 == 0 {
				rn := rand.Intn(100)
				newdata := []byte{}
				for i := 0; i < rn; i++ {
					r := rand.Intn(40)
					newdata = append(newdata, byte(r))
				}
				node.Data = newdata
			}
			node.PingCount++

		}
	}
}

func (node *Node) checkForPing() {

	//while true
	for true {

		//if its master, break
		if node.Rank == MASTER {
			return
		}

		//reset pinged value and wait for 2 seconds
		// node.Pinged = false
		// time.Sleep(2 * time.Second)
		// fmt.Println(NODE_MAP)
		//if the ping hasnt changed (eg ping didnt set it to true)
		// fmt.Println(time.Now().Sub(node.Pinged).Seconds())
		if time.Now().Sub(node.Pinged).Seconds() > 2.5 {

			//node must be down
			fmt.Println("NO PING FROM MASTER, MUST BE DOWN!!!", node.Ip, NODE_MAP[1].Ip)

			//if node is next in line to throne
			if node.Ip == NODE_MAP[1].Ip {

				//set that node to master
				node.Rank = MASTER
				node.Data = NODE_MAP[0].Data

				fmt.Println("IM THE MASTER NOW, COPIED ALL DATA FROM PREVIOUS MASTER!!!")

			} else {
				node.Pinged = time.Now()
			}

			//update node map
			NODE_MAP = NODE_MAP[1:]
			NODE_MAP[0] = node
			fmt.Println("UPDATING NODE MAP", getNodeIps())
			//wait for a second to keep time
			// time.Sleep(1 * time.Second)
		}

	}
}

func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {

	//only receive ping if its a follower
	if node.Rank == FOLLOWER {

		//fetch data of ping (nodemap)
		body, _ := ioutil.ReadAll(req.Body)

		//move nodemap to local memory
		var localnm []*Node
		json.Unmarshal(body, &localnm)

		if len(localnm[0].Data) == 0 {

			// fmt.Println("data is empty, must not have changed")
			node.Pinged = time.Now()
			// fmt.Println("PINGED")
			w.Header().Add("pinged", "true")
			return
		}
		fmt.Println("UPDATED DATA ON THIS NODE", localnm[0].Data)
		// fmt.Println("data mustve changed", localnm[0].Data)
		NODE_MAP = []*Node{}
		for _, j := range localnm {
			NODE_MAP = append(NODE_MAP, j)
		}

		//currentmaster is top of list
		currentmasterip = NODE_MAP[0].Ip

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
		fmt.Println(ip, " has connected")
		NODE_MAP = append(NODE_MAP, &Node{ip, time.Now(), 0, FOLLOWER, []byte{}})
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
		fmt.Println("RECEIVED DATA UPDATE: ", req.Header.Get("data"))
	}
}

func main() {

	//init node and set to master (true until proven otherwise)
	var node Node
	node.Rank = FOLLOWER
	connectoip := ""
	masterip := ""

	fmt.Println("Type in ip to connect, or enter if master.")
	fmt.Scanf("%s", &connectoip)
	if connectoip == "" {
		//get ip of machine
		connectoip = "localhost"
		masterip = "localhost:3100"
	} else {
		resp, err := http.Get("http://" + connectoip + "/connect")

		if err == nil && resp.Header.Get("rank") == "MASTER" {

			masterip = connectoip

		} else {
			fmt.Println("THAT IS NOT A MASTER")
			return
		}
	}
	go node.pickPort("localhost")

	// //check for current masters
	// for i := 0; i < 3; i++ {

	// 	//check whether is port is already open
	// 	con, err := net.DialTimeout("tcp", connectoip+":"+strconv.Itoa(PORT_START+i), 100*time.Millisecond)

	// 	//if it is, continue
	// 	if err != nil {
	// 		continue
	// 	}

	// 	//check whether the open port is a master
	// 	resp, err := http.Get("http://localhost:" + strconv.Itoa(PORT_START+i) + "/connect")

	// 	if err == nil && resp.Header.Get("rank") == "MASTER" {

	// 		masterip = masterip + ":" + strconv.Itoa(PORT_START+i)

	// 		break
	// 	}

	// 	//close each connection (only for checking)
	// 	con.Close()

	// }
	fmt.Println("The master is ", masterip)

	//pick a port for this node to be on

	//this needs to be changed, basically, wait for port selection to be done
	time.Sleep(3 * time.Second)
	//add to local nodemap (will be replicated if its master)
	NODE_MAP = append(NODE_MAP, &node)

	//if this node is not on the master port, then its a follower
	if masterip == node.Ip {
		node.Rank = MASTER
		node.Data = []byte{0x64, 0x43, 0x10}
	}

	//connect to port
	if node.Rank == FOLLOWER {
		fmt.Println("TRYING TO CONNECT TO ", "http://"+masterip)
		_, err := http.Get("http://" + *&masterip + "/connect?ip=" + node.Ip)
		if err != nil {
			fmt.Println(err)
		}
	}

	//http handlers
	http.HandleFunc("/connect", node.connectHandler)
	http.HandleFunc("/ping", node.pingHandler)
	http.HandleFunc("/data", node.dataHandler)
	node.Pinged = time.Now()
	//ping handling
	go node.ping()
	go node.checkForPing()

	for true {
	}
}

//TODO

// NEED TO WAIT FOR MASTER TO COME ONLINE
// GET DATA FROM MASTER (for some api or sdk)
