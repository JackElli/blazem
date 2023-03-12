package main

import (
	"blazem/endpoints"
	"blazem/global"
	"blazem/logging"
	"blazem/query"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Node global.Node

// utils
// return the index of the node in the nodemap
func indexOfNodeInNodeMap(node *global.Node) int {
	for i, n := range global.NODE_MAP {
		if n.Ip == node.Ip {
			return i
		}
	}
	return -1
}

// return the data stored in the nodemap
func getNodeDatas() []sync.Map {
	var nodedata []sync.Map
	for _, n := range global.NODE_MAP {
		nodedata = append(nodedata, n.Data)
	}
	return nodedata
}

func getLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// ping connections
func (node *Node) tryListen(ip string) {

	// //wait for 500 milliseconds
	// time.Sleep(500 * time.Millisecond)

	//listen on selected port
	portstr := ip
	if strings.Count(ip, ":") > 1 {
		portstr = strings.Split(ip, ":")[0]
	}

	global.Logger.Log("trying on "+portstr, logging.INFO)
	l, err := net.Listen("tcp", portstr)

	//if theres an error in connecting, stop
	if err != nil {
		// defer l.Close()
		return
	}

	//if there's no error set the nodes port to the current port
	node.Ip = ip
	global.Logger.Log("Blazem started up on "+ip, logging.INFO)

	//serve http requests on this port
	http.Serve(l, nil)
}

func (node *Node) pickPort(ip string) {

	connectIp := ""
	//for each node, try to listen
	for i := 0; i < 3; i++ {
		connectIp = ip + ":" + strconv.Itoa(global.PORT_START+i)
		node.tryListen(connectIp)

		//if theres an error, break out
		if node.Ip != "" {
			break
		}
	}
}

func (node *Node) setNodeMasterAttrs() {
	node.Rank = global.MASTER
	node.Data = sync.Map{}

	testData1 := map[string]interface{}{
		"type":       "folder",
		"key":        "testkey1",
		"folderName": "TestFolder",
		"value":      "hello this is a test",
		"date":       time.Now().Format("2006-01-02T15:04:05"),
	}

	testData2 := map[string]interface{}{
		"type":   "text",
		"key":    "testkey2",
		"folder": "TestFolder",
		"value":  "hello asd",
		"date":   time.Now().Format("2006-01-02T15:04:05"),
	}
	testData3 := map[string]interface{}{
		"type":       "folder",
		"key":        "testkey3",
		"folderName": "AnotherFolder",
		"folder":     "TestFolder",
		"date":       time.Now().Format("2006-01-02T15:04:05"),
	}
	testData4 := map[string]interface{}{
		"type":   "text",
		"key":    "testkey4",
		"folder": "AnotherFolder",
		"value":  "hello testing",
		"date":   time.Now().Format("2006-01-02T15:04:05"),
	}
	node.Data.Store("testkey", testData1)
	node.Data.Store("testkey2", testData2)
	node.Data.Store("testkey3", testData3)
	node.Data.Store("testkey4", testData4)
}
func setupLogger() {
	//setup file for logging
	logfile := "logging/"
	global.Logger = *logging.LogFile(logfile)
}

// main func
func main() {

	//init node and set to follower (true until proven otherwise)
	var node Node = Node{
		Ip:            "",
		Pinged:        time.Now(),
		PingCount:     0,
		Rank:          global.FOLLOWER,
		Data:          sync.Map{},
		Active:        true,
		RecentQueries: map[string]string{},
		Rules:         map[string]global.Rule{},
	}

	global.GlobalNode = (*global.Node)(&node)

	// setup the logger
	setupLogger()

	var masterip string = ""
	//set ips
	localip := getLocalIp()

	//this needs to be async as port should be on other thread
	go node.pickPort(localip)

	//setup endpoints
	endpoints.SetupHandlers((*endpoints.Node)(&node))

	//add to local nodemap (will be replicated if its master)
	global.NODE_MAP = append(global.NODE_MAP, (*global.Node)(&node))

	//if this node is not on the master port, then its a follower
	if masterip == node.Ip {
		node.setNodeMasterAttrs()
	}

	// load local storage
	(*global.Node)(&node).ReadFromLocal()

	//ping handling
	go (*global.Node)(&node).Ping()

	// run the rule checker
	go (*endpoints.Node)(&node).CheckRules()

	query.LoadIntoMemory(global.Node(node))

	//like a game loop
	for true {
		time.Sleep(150 * time.Millisecond)
	}
}
