package main

import (
	"distributed_servers/endpoints"
	"distributed_servers/global"
	"distributed_servers/logging"
	"distributed_servers/webend"
	"fmt"
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
		newmap = append(newmap, global.Node{Ip: n.Ip, Pinged: n.Pinged, PingCount: n.PingCount, Rank: n.Rank, Data: n.Data, Active: true})
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

	//wait for 500 milliseconds
	time.Sleep(500 * time.Millisecond)

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
	global.Logger.Log("connected on "+portstr, logging.GOOD)

	//serve http requests on this port
	http.Serve(l, nil)
}

func (node *Node) pickPort(ip string) {

	//for each node, try to listen
	for i := 0; i < 3; i++ {
		node.tryListen(ip + ":" + strconv.Itoa(global.PORT_START+i))

		//if theres an error, break out
		if node.Ip != "" {
			break
		}
	}
}

func (node *Node) connectToIp(masterip *string, localip string) {
	var connectoip string = ""
	fmt.Println("Type in ip to connect, or enter if master.")
	fmt.Scanf("%s", &connectoip)

	if connectoip == "" {
		//get ip of machine
		//THIS IS THE IP TO CONNECT TO
		connectoip = localip + ":3100"
		*masterip = localip + ":3100"
	} else {
		//connect to server
		global.Logger.Log("TRYING TO CONNECT TO "+"http://"+connectoip, logging.INFO)
		resp, err := http.Get("http://" + connectoip + "/connect?ip=" + node.Ip)

		if err == nil && resp.Header.Get("rank") == "MASTER" {
			*masterip = connectoip
		} else {
			global.Logger.Log("THAT IS NOT A MASTER", logging.ERROR)
			return
		}
		defer resp.Body.Close()
	}
}

func (node *Node) setNodeMasterAttrs() {
	node.Rank = global.MASTER
	node.Data = global.NodeData{}

	testData1 := global.JsonData{
		Folder: "text",
		Data:   "hello this is a test",
		Type:   "text",
		Date:   time.Now(),
	}

	testData2 := global.JsonData{
		Folder: "text",
		Data:   "hello asd",
		Type:   "text",
		Date:   time.Now(),
	}

	node.Data["testkey"] = testData1
	node.Data["testkey2"] = testData2
}

// main func
func main() {

	//setup file for logging
	logfile := "logging/"
	global.Logger = *logging.LogFile(logfile)

	//init node and set to follower (true until proven otherwise)
	var node Node

	var masterip string = ""

	node.Rank = global.FOLLOWER
	node.Active = true
	//set ips
	localip := getLocalIp()
	//this needs to be async as port should be on other thread
	go node.pickPort(localip)

	time.Sleep(2 * time.Second)
	node.connectToIp(&masterip, localip)
	global.Logger.Log("The master is "+masterip, logging.INFO)

	//setup endpoints
	endpoints.SetupHandlers((*endpoints.Node)(&node))
	webend.SetupWebend()

	//add to local nodemap (will be replicated if its master)
	global.NODE_MAP = append(global.NODE_MAP, (*global.Node)(&node))

	//if this node is not on the master port, then its a follower
	if masterip == node.Ip {
		node.setNodeMasterAttrs()
	}

	//ping handling
	node.Pinged = time.Now()
	go (*global.Node)(&node).Ping()

	//like a game loop
	for true {
		time.Sleep(150 * time.Millisecond)
	}
}

//TODO

// NEED TO WAIT FOR MASTER TO COME ONLINE
// GET DATA FROM MASTER (for some api or sdk)

//BUG WHY NEW MASTER WHEN NODE NOT DOWN???

//NEED TO DO DISK PERSISTENCE

// DISTRIBUTED COMPUTING, ADDITION OF NUMBERS ON 3 CPUS??
// CLEAN CODE

//MIGHT WANT TO CHANGE IT SO WE CHECK EVERY MINUTE RATHER THAN 4 SECONDS

//THERES A BUG SOMEWHERE WITH NULLPOINTER

//SWITCH FROM FDS (FULL DATA SYNC) TO DDS (DELTA DATA SYNC)

//CHECK NODE IS NOT ON TWO CLUSTERS

//change code to allocate mem before array swap
