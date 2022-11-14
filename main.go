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
	"text/template"
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

var logger logging.Logger

var dataChanged bool = false

type Node struct {
	Ip        string
	Pinged    time.Time
	PingCount int
	Rank      Rank
	Data      map[string]string
	Active    bool
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
		newmap = append(newmap, &Node{n.Ip, n.Pinged, 0, n.Rank, map[string]string{}, true})
	}
	return newmap
}

func nodeMapPointertoMem() []Node {
	var newmap []Node
	for _, n := range NODE_MAP {
		newmap = append(newmap, Node{n.Ip, n.Pinged, n.PingCount, n.Rank, n.Data, true})
	}
	return newmap
}

func alreadyInNodeMap(ip string) bool {
	for _, n := range NODE_MAP {
		if n.Ip == ip {
			return true
		}
	}
	return false
}

func getAllDataToPrint(data map[string]string) []string {
	var retdata []string
	for v := range data {
		retdata = append(retdata, data[v])
	}
	return retdata
}

func getLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

func checkIfDataChanged() []byte {
	var jsonNodeMap []byte
	if dataChanged {
		logger.Log("DATA CHANGED", logging.INFO)
		jsonNodeMap, _ = json.Marshal(NODE_MAP)
		dataChanged = false
	} else {
		jsonNodeMap, _ = json.Marshal(getNodeMapWithoutData())
	}
	return jsonNodeMap
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

	logger.Log("trying on "+portstr, logging.INFO)
	l, err := net.Listen("tcp", portstr)

	//if theres an error in connecting, stop
	if err != nil {
		// defer l.Close()
		return
	}

	//if there's no error set the nodes port to the current port
	node.Ip = ip
	logger.Log("connected on "+portstr, logging.GOOD)

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

func (n *Node) pingRetry(resp *http.Response, c *http.Client, sendData *bytes.Buffer) *http.Response {
	logger.Log("PINGING again "+"http://"+n.Ip, logging.INFO)
	fmt.Println(c, sendData)
	resp, err := c.Post("http://"+n.Ip+"/ping", "application/json", sendData)
	// resp = resps
	// fmt.Println(resps.Header.Get("pinged"))
	if err != nil {
		indexOfNode := indexOfNodeInNodeMap(n)
		NODE_MAP = append(NODE_MAP[:indexOfNode], NODE_MAP[indexOfNode+1:]...)
		return nil
	}
	return resp
}

func (n *Node) receivedPing(resp *http.Response, c *http.Client) {
	//need to do master resilience here
	if resp.Header.Get("pinged") == "true" {
		logger.Log("PING RECEIVED FROM "+n.Ip, logging.GOOD)
		c.CloseIdleConnections()
	}
}

func (node *Node) pingEachConnection(jsonNodeMap []byte) {
	for _, n := range NODE_MAP {

		//don't ping to itself
		if n.Ip == node.Ip || n.Active == false {
			continue
		}
		//marshall so we're able to send over TCP
		if n.PingCount == 0 {
			logger.Log("SENDING MAP TO FIRST JOINER", logging.INFO)
			jsonNodeMap, _ = json.Marshal(NODE_MAP)
		}
		sendData := bytes.NewBuffer(jsonNodeMap)
		//ping connection
		logger.Log("PINGING "+n.Ip, logging.INFO)
		_, err := net.DialTimeout("tcp", n.Ip, 1500*time.Millisecond)
		//increase connection ping count
		n.PingCount++
		//retry logic
		if err != nil {
			_, err := net.DialTimeout("tcp", n.Ip, 2000*time.Millisecond)
			if err != nil {
				logger.Log("Cannot connect to "+n.Ip, logging.WARNING)
				//dont remove node, just make it inactive
				// indexOfNode := indexOfNodeInNodeMap(n)
				// NODE_MAP = append(NODE_MAP[:indexOfNode], NODE_MAP[indexOfNode+1:]...)
				n.Active = false
				continue
			}
		}
		logger.Log("PING RECEIVED FROM "+n.Ip, logging.INFO)
		_, err = http.Post("http://"+n.Ip+"/ping", "application/json", sendData)
	}
	node.PingCount++
}

func (node *Node) ping() {

	//while true
	for true {

		time.Sleep(4 * time.Second)
		//print the rank of the node and wait for 2 secs
		if node.Rank == MASTER {
			logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+strings.Join(getNodeIps(), " "), logging.INFO)
		}
		if len(NODE_MAP) == 1 {
			continue
		}
		if node.Rank == FOLLOWER {
			continue
		}
		// fmt.Println("Data stored: ", node.Data)
		//check if the data has changed from the data on the map
		jsonNodeMap := checkIfDataChanged()
		//need to check numbers
		node.pingEachConnection(jsonNodeMap)
	}
}

func (node *Node) checkForNoPingFromMaster() {

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
	logger.Log("Slow response first check at "+fmt.Sprintf("%f", timeSinceLastPingAbs)+"s", logging.WARNING)
	//if not, check for retry ping
	time.Sleep(4 * time.Second)
	timeSinceLastPingAbs = time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 8 {
		return
	}
	//if no pings in that time, master is down
	logger.Log("NO PING FROM MASTER!!!", logging.INFO)

	//if node is not next in line, break out
	if node.Ip != NODE_MAP[1].Ip {
		return
	}
	//set that node to master
	node.Rank = MASTER
	node.Data = NODE_MAP[0].Data
	waitingTimeStr := strconv.Itoa(int(time.Now().Sub(node.Pinged).Seconds()))
	logger.Log("IM THE MASTER NOW, COPIED ALL DATA FROM PREVIOUS MASTER!!! after waiting for "+waitingTimeStr+"s", logging.GOOD)
	//update node map
	NODE_MAP = NODE_MAP[1:]
	NODE_MAP[0] = node
}

// handlers
func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {

	//only receive ping if its a follower
	if node.Rank == FOLLOWER {
		//print the rank of the node and wait for 2 secs
		logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+strings.Join(getNodeIps(), " "), logging.INFO)
		node.Pinged = time.Now()
		//need to check for ping here (start )
		go node.checkForNoPingFromMaster()
		//fetch data of ping (nodemap)
		body, _ := ioutil.ReadAll(req.Body)

		//move nodemap to local memory
		var localnm []*Node
		json.Unmarshal(body, &localnm)

		//add the changed node map
		currentMasterData := NODE_MAP[0].Data
		NODE_MAP = localnm

		if len(localnm[0].Data) == 0 {
			NODE_MAP[0].Data = currentMasterData
			// fmt.Println(getNodeDatas())
			// fmt.Println("data is empty, must not have changed")

			// fmt.Println("PINGED")
			w.Header().Add("pinged", "true")
			// fmt.Println("sent ping back")
			return
		}

		logger.Log("UPDATED DATA ON THIS NODE!", logging.GOOD)

		NODE_MAP = []*Node{}
		for _, j := range localnm {
			NODE_MAP = append(NODE_MAP, j)
		}
		//it has been successfully pinged
		w.Header().Add("pinged", "true")
		// fmt.Println("sent ping back")
	} else {
		fmt.Println("SOMETHINGS GONE WRONG!")
		node.Rank = FOLLOWER
	}

	//NEED TO ENFORCE MASTER RESILIENCE (if two masters exist, pick one)
}

func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {

	//find the port of client trying to connect
	ip := req.URL.Query().Get("ip")

	//if port isnt 0 (there is a valid port)
	if ip != "" {
		//add to the node map
		logger.Log(ip+" has connected", logging.GOOD)
		if !alreadyInNodeMap(ip) {
			NODE_MAP = append(NODE_MAP, &Node{ip, time.Now(), 0, FOLLOWER, map[string]string{}, true})
		} else {
			//already in map
			indexOfNode := indexOfNodeInNodeMap(&Node{ip, time.Now(), 0, FOLLOWER, nil, true})
			NODE_MAP[indexOfNode].Active = true
			//set to active
		}
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

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, all")

		if req.Method == "POST" {

			//TODO send multiple key and values

			var dataToSet []string
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &dataToSet)

			setKey := dataToSet[0]
			setVal := dataToSet[1]
			dataType := dataToSet[2]

			//ADD DATE OF ENTRY
			value := "{\"data\":\"" + setVal + "\", \"type\":\"" + dataType + "\", \"date\":\"" + time.Now().String() + "\"}"
			node.Data[setKey] = value

			dataChanged = true
			w.Header().Set("response", "done")
			return
		}

		datagetkey := req.Header.Get("key")
		datagetall := req.Header.Get("all")

		if datagetkey != "" {
			dataToSend := node.Data[datagetkey]
			w.Header().Set("val", dataToSend)
		}
		if datagetall != "" {
			dataToSend := getAllDataToPrint(node.Data)
			// w.Header().Set("val", strings.Join(dataToSend, ","))
			json.NewEncoder(w).Encode(dataToSend)
		}

	}
}

func nodeMapHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, all")

	nodeMapResp := getNodeIps()

	json.NewEncoder(w).Encode(nodeMapResp)
}

func webHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("statictest/index.html")

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// main func
func main() {

	//setup file for logging
	logfile := "logging/"
	logger = *logging.LogFile(logfile)

	//init node and set to follower (true until proven otherwise)
	var node Node
	node.Rank = FOLLOWER

	//set ips
	localip := getLocalIp()
	connectoip := ""
	masterip := ""
	go node.pickPort(localip)
	time.Sleep(2 * time.Second)
	fmt.Println("Type in ip to connect, or enter if master.")
	fmt.Scanf("%s", &connectoip)
	if connectoip == "" {
		//get ip of machine
		//THIS IS THE IP TO CONNECT TO
		connectoip = localip + ":3100"
		masterip = localip + ":3100"
	} else {
		//connect to server
		logger.Log("TRYING TO CONNECT TO "+"http://"+connectoip, logging.INFO)
		resp, err := http.Get("http://" + connectoip + "/connect?ip=" + node.Ip)
		if err == nil && resp.Header.Get("rank") == "MASTER" {
			masterip = connectoip
			// fmt.Println("Yes")
		} else {
			logger.Log("THAT IS NOT A MASTER", logging.ERROR)
			return
		}
	}

	//pick a port for this node to be on you local ip
	logger.Log("The master is "+masterip, logging.INFO)

	//handle incoming connection (MAYBE MAKE THIS ASYNC)
	http.HandleFunc("/connect", node.connectHandler)
	http.HandleFunc("/ping", node.pingHandler)
	http.HandleFunc("/data", node.dataHandler)
	http.HandleFunc("/nodemap", nodeMapHandler)

	//web stuff
	styles := http.FileServer(http.Dir("./statictest/css"))
	scripts := http.FileServer(http.Dir("./statictest/js"))
	http.HandleFunc("/", webHandler)
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))

	//this needs to be changed, basically, wait for port selection to be done
	// time.Sleep(3 * time.Second)
	//add to local nodemap (will be replicated if its master)
	NODE_MAP = append(NODE_MAP, &node)

	//if this node is not on the master port, then its a follower
	if masterip == node.Ip {
		node.Rank = MASTER
		node.Data = map[string]string{}
		node.Data["testkey"] = "{\"data\":\"hello this is a test\", \"type\":\"text\", \"date\":\"" + time.Now().String() + "\"}"
		node.Data["testkey2"] = "{\"data\":\"hello asd\", \"type\":\"text\", \"date\":\"" + time.Now().String() + "\"}"
	}

	//ping handling
	node.Pinged = time.Now()

	go node.ping()

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
