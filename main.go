package main

import (
	"distributed_servers/endpoints"
	"distributed_servers/global"
	"distributed_servers/logging"
	"distributed_servers/webend"
	"net"
	"net/http"
	"strconv"
	"strings"
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
func getNodeDatas() []global.NodeData {
	var nodedata []global.NodeData
	for _, n := range global.NODE_MAP {
		nodedata = append(nodedata, n.Data)
	}
	return nodedata
}

func nodeMapPointertoMem() []global.Node {
	var newmap []global.Node
	for _, n := range global.NODE_MAP {
		newmap = append(newmap, global.Node{Ip: n.Ip, Pinged: n.Pinged,
			PingCount: n.PingCount, Rank: n.Rank, Data: n.Data, Active: true})
	}
	return newmap
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
	node.Data = global.NodeData{}

	testData1 := map[string]interface{}{
		"key":    "testkey",
		"folder": "text",
		"value":  "hello this is a test",
		"type":   "text",
		"date":   time.Now(),
	}

	testData2 := map[string]interface{}{
		"key":    "testkey2",
		"folder": "text",
		"value":  "hello asd",
		"type":   "text",
		"date":   time.Now(),
	}

	node.Data["testkey"] = testData1
	node.Data["testkey2"] = testData2
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
		Data:          map[string]interface{}{},
		Active:        true,
		RecentQueries: map[string]time.Time{},
		Rules:         map[string][]global.Task{},
	}

	// setup the logger
	setupLogger()

	var masterip string = ""
	//set ips
	localip := getLocalIp()

	//this needs to be async as port should be on other thread
	go node.pickPort(localip)

	//setup endpoints
	endpoints.SetupHandlers((*endpoints.Node)(&node))
	webend.SetupWebend()

	//add to local nodemap (will be replicated if its master)
	global.NODE_MAP = append(global.NODE_MAP, (*global.Node)(&node))

	//if this node is not on the master port, then its a follower
	if masterip == node.Ip {
		node.setNodeMasterAttrs()
	}

	// create backup
	(*global.Node)(&node).SaveDataJson()

	//ping handling
	go (*global.Node)(&node).Ping()

	//like a game loop
	for true {
		time.Sleep(150 * time.Millisecond)
	}
}
