package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Rank string

const (
	MASTER   Rank = "MASTER"
	FOLLOWER Rank = "FOLLOWER"
)

var PORT_START = 3100
var NODE_MAP []*Node

var currentmasterport int = 3100

type Node struct {
	Port   *int
	Pinged bool
	Rank   Rank
	Data   []byte
}

// need to connect to master
func (node *Node) tryListen(port int) {

	//wait for 500 milliseconds
	time.Sleep(500 * time.Millisecond)

	//listen on selected port
	portstr := ":" + strconv.Itoa(port)
	fmt.Println("trying on " + portstr)
	l, err := net.Listen("tcp", portstr)

	//if theres an error in connecting, stop
	if err != nil {
		return
	}

	//if there's no error set the nodes port to the current port
	node.Port = &port
	fmt.Println("connected on " + portstr)

	//serve http requests on this port
	http.Serve(l, nil)
}

func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {

	//find the port of client trying to connect
	connectionPort := req.URL.Query().Get("port")
	portInt, _ := strconv.Atoi(connectionPort)

	//if port isnt 0 (there is a valid port)
	if portInt != 0 {

		//add to the node map
		fmt.Println(connectionPort, " has connected")
		NODE_MAP = append(NODE_MAP, &Node{&portInt, false, FOLLOWER, []byte{}})
	}

	//write back to the client with the rank of the node its trying to connect to
	w.Header().Add("rank", string(node.Rank))
}

func (node *Node) pickPort() {

	//for each node, try to listen
	for i := 0; i < 3; i++ {
		node.tryListen(PORT_START + i)

		//if theres an error, break out
		if node.Port != nil {
			break
		}
	}
}

func indexOfNodeInNodeMap(node *Node) int {
	for i, n := range NODE_MAP {
		if *n.Port == *node.Port {
			return i
		}
	}
	return -1
}

func getNodeDatas() [][]byte {
	var nodedata [][]byte
	for _, n := range NODE_MAP {
		nodedata = append(nodedata, n.Data)
	}
	return nodedata
}
func (node *Node) ping() {

	//while true
	for true {

		//print the rank of the node and wait for 2 secs
		// fmt.Println(node.Rank, getNodeDatas())
		fmt.Println("This node has data: ", node.Data)
		time.Sleep(2 * time.Second)

		//only send ping if its a master
		if node.Rank == MASTER {
			// fmt.Println("SENDING PING")

			//each port in the node map
			for _, n := range NODE_MAP {

				//marshall so we're able to send over TCP
				jsonNodeMap, _ := json.Marshal(NODE_MAP)
				sendData := bytes.NewBuffer(jsonNodeMap)

				//don't ping to itself
				if *n.Port != *node.Port {
					resp, err := http.Post("http://localhost:"+strconv.Itoa(*n.Port)+"/ping", "application/json", sendData)

					if err != nil {
						indexOfNode := indexOfNodeInNodeMap(n)
						NODE_MAP = append(NODE_MAP[:indexOfNode], NODE_MAP[indexOfNode+1:]...)
						continue
					}
					if resp.Header.Get("pinged") == "true" {
						fmt.Println("PING RECEIVED FROM ", *n.Port)
					}
				}
			}
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
		node.Pinged = false
		time.Sleep(2 * time.Second)

		//if the ping hasnt changed (eg ping didnt set it to true)
		if node.Pinged == false {

			//node must be down
			fmt.Println("NO PING FROM MASTER, MUST BE DOWN!!!")

			//if node is next in line to throne
			if *node.Port == *NODE_MAP[1].Port {

				//set that node to master
				node.Rank = MASTER
				node.Data = NODE_MAP[0].Data

				fmt.Println("IM THE MASTER NOW!!!")
			}

			//prune node map
			NODE_MAP = NODE_MAP[1:]

			//add data to master (top of list)
			NODE_MAP[0].Data = node.Data

			//wait for a second to keep time
			time.Sleep(1 * time.Second)
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
		NODE_MAP = []*Node{}
		for _, j := range localnm {
			NODE_MAP = append(NODE_MAP, j)
		}

		//currentmaster is top of list
		currentmasterport = *NODE_MAP[0].Port

		// the node map is ", NODE_MAP
		//print to user
		// fmt.Println("Received ping from master {", currentmasterport, "}")

		//it has been successfully pinged
		node.Pinged = true
		w.Header().Add("pinged", "true")
	}

}
func main() {

	//init node and set to master (true until proven otherwise)
	var node Node
	node.Rank = MASTER
	masterport := -1

	//check for current masters
	for i := 0; i < 3; i++ {

		//check whether is port is already open
		con, err := net.DialTimeout("tcp", ":"+strconv.Itoa(PORT_START+i), 100*time.Millisecond)

		//if it is, continue
		if err != nil {
			continue
		}

		//check whether the open port is a master
		resp, err := http.Get("http://localhost:" + strconv.Itoa(PORT_START+i) + "/connect")
		if err == nil && resp.Header.Get("rank") == "MASTER" {
			masterport = PORT_START + i
			break
		}

		//close each connection (only for checking)
		con.Close()

	}

	if masterport == -1 {
		node.Data = []byte{0x64, 0x43, 0x10}
		masterport = 3100
	}
	fmt.Println("The master is ", masterport)

	//pick a port for this node to be on
	go node.pickPort()

	//this needs to be changed, basically, wait for port selection to be done
	time.Sleep(3 * time.Second)
	//add to local nodemap (will be replicated if its master)
	NODE_MAP = append(NODE_MAP, &node)

	//if this node is not on the master port, then its a follower
	if masterport != *node.Port {
		node.Rank = FOLLOWER
	}

	//connect to port
	_, err := http.Get("http://localhost:" + strconv.Itoa(masterport) + "/connect?port=" + strconv.Itoa(*node.Port))
	if err != nil {
		fmt.Println(err)
	}

	//http handlers
	http.HandleFunc("/connect", node.connectHandler)
	http.HandleFunc("/ping", node.pingHandler)

	//ping handling
	go node.ping()
	go node.checkForPing()

	for true {
	}
}
