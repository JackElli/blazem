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
var NODE_MAP []int

var currentmasterport int = 3100

type Node struct {
	Port   *int
	Master *int
	Rank   Rank
	Pinged bool
}

// need to connect to master
func (node *Node) tryListen(port int) {
	time.Sleep(500 * time.Millisecond)
	portstr := ":" + strconv.Itoa(port)
	fmt.Println("trying on " + portstr)
	l, err := net.Listen("tcp", portstr)
	if err != nil {
		// node.Rank = FOLLOWER
		// node.Master = &port
		return
	}
	node.Port = &port
	fmt.Println("connected on " + portstr)
	http.Serve(l, nil)
}

func (node *Node) checkForMaster() {
	for true {
		time.Sleep(2 * time.Second)
		if node.Master != nil {
			fmt.Println("My master is " + strconv.Itoa(*node.Master))
		} else {
			fmt.Println("I AM MASTER")
		}
	}
}

//	func (node *Node) ping() {
//		fmt.Print("My master is " + strconv.Itoa(*node.Master))
//	}
func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {
	connectionPort := req.URL.Query().Get("port")
	portInt, _ := strconv.Atoi(connectionPort)
	if portInt != 0 {
		fmt.Println(connectionPort, " is trying to connect")
		NODE_MAP = append(NODE_MAP, portInt)
		fmt.Println(NODE_MAP)

	}
	w.Header().Add("rank", string(node.Rank))
}

func (node *Node) pickPort() {
	for i := 0; i < 3; i++ {
		node.tryListen(PORT_START + i)

		if node.Port != nil {
			break
		}
	}
}

type JsonNode struct {
	Port   *int
	Master *int
	Rank   *Rank
	Pinged bool
}

func (node *Node) ping() {

	for true {
		fmt.Println(node.Rank)
		time.Sleep(2 * time.Second)

		if node.Rank == MASTER {
			fmt.Println("SENDING PING")
			for _, n := range NODE_MAP {

				jsonNodeMap, _ := json.Marshal(NODE_MAP)

				// fmt.Println(jsonNodeMap)
				sendData := bytes.NewBuffer(jsonNodeMap)
				if n != *node.Port {
					_, _ = http.Post("http://localhost:"+strconv.Itoa(n)+"/ping", "application/json", sendData)
				}
			}
		}
	}
}

func (node *Node) checkForPing() {

	for true {
		// fmt.Println("PING RANK", node.Rank)
		if node.Rank == MASTER {
			return
		}
		node.Pinged = false
		time.Sleep(2 * time.Second)
		if node.Pinged == false {
			fmt.Println("NO PING FROM MASTER, MUST BE DOWN!!!", NODE_MAP)
			fmt.Println(*node.Port, NODE_MAP[1])
			if *node.Port == NODE_MAP[1] {
				fmt.Println(*node.Port, " ", NODE_MAP[1])
				node.Rank = MASTER

			}
			NODE_MAP = NODE_MAP[1:]
			time.Sleep(1 * time.Second)
		}
	}
}

func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {
	if node.Rank == FOLLOWER {
		body, _ := ioutil.ReadAll(req.Body)
		// fmt.Println(body)
		var localnm []int
		json.Unmarshal(body, &localnm)
		NODE_MAP = []int{}
		for _, j := range localnm {
			NODE_MAP = append(NODE_MAP, j)
		}
		currentmasterport = NODE_MAP[0]
		// NODE_MAP = nodemap

		fmt.Println("Received ping from master!! ", currentmasterport, " the node map is ", NODE_MAP)

		node.Pinged = true
	}

}
func main() {
	var node Node
	node.Rank = MASTER
	// node.NodeMap = append(node.NodeMap, &node)

	masterport := 3100
	//check for current masters
	for i := 0; i < 3; i++ {

		// _, err := http.Get("http://localhost:" + strconv.Itoa(PORT_START+i) + "/connect")
		// fmt.Println("tested port ", PORT_START+i, " got err: ", err)
		con, err := net.DialTimeout("tcp", ":"+strconv.Itoa(PORT_START+i), 100*time.Millisecond)
		// fmt.Println("con: ", con, " err: ", err)
		fmt.Println("PORT ", PORT_START+i, " has err: ", err)
		if err != nil {
			// con.Close()
			continue
		}
		resp, err := http.Get("http://localhost:" + strconv.Itoa(PORT_START+i) + "/connect")
		if err == nil && resp.Header.Get("rank") == "MASTER" {
			masterport = PORT_START + i
			break
		}
		con.Close()

	}
	fmt.Println("The master is ", masterport)

	go node.pickPort()
	time.Sleep(3 * time.Second)
	NODE_MAP = append(NODE_MAP, *node.Port)
	if masterport != *node.Port {
		node.Rank = FOLLOWER
	}
	resp, err := http.Get("http://localhost:" + strconv.Itoa(masterport) + "/connect?port=" + strconv.Itoa(*node.Port))
	fmt.Println("RESPONSE: ", resp.Header.Get("rank"), err)
	http.HandleFunc("/connect", node.connectHandler)
	http.HandleFunc("/ping", node.pingHandler)

	go node.ping()
	go node.checkForPing()
	// if node.Rank == MASTER {
	// 	go node.ping()
	// }
	// fmt.Println("Rank" + node.Rank)
	// if node.Rank == FOLLOWER {

	// }

	// go tryListen("asdsad")

	for true {
	}
}
