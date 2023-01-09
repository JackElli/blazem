package endpoints

import (
	"bytes"
	"distributed_servers/global"
	"distributed_servers/logging"
	"distributed_servers/query"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Node global.Node

type WebNodeMap struct {
	Ip     string `json:"ip"`
	Active bool   `json:"active"`
}

type SendData struct {
	Key  string          `json:"key"`
	Data global.JsonData `json:"data"`
}

type Stats struct {
	Cpu float64 `json:"cpu"`
	Ram float64 `json:"ram"`
}

type SendQueryData struct {
	Docs      []SendData `json:"docs"`
	TimeTaken int64      `json:"timeTaken"`
}

func getHexKey() string {
	pos := "0123456789abcdef"
	key := ""
	for i := 0; i < 16; i++ {
		key += string(pos[rand.Intn(len(pos)-1)])
	}
	return key
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func getWindowsStats() Stats {

	ps, _ := exec.LookPath("powershell.exe")
	cpu := exec.Command(ps, "Get-CimInstance win32_processor | Measure-Object -Property LoadPercentage -Average")
	ramTotal := exec.Command(ps, "wmic ComputerSystem get TotalPhysicalMemory")
	ramFree := exec.Command(ps, "wmic OS get FreePhysicalMemory")

	//CPU
	var cpuout bytes.Buffer
	cpu.Stdout = &cpuout
	cpuerr := cpu.Run()
	if cpuerr != nil {
		fmt.Println(cpuerr)
	}

	cpuavreg, _ := regexp.Compile("Average  : [0-9]*")
	cpuav := cpuavreg.FindString(cpuout.String())
	cpureg, _ := regexp.Compile("[0-9]+")
	cpuStat, _ := strconv.ParseFloat(cpureg.FindString(cpuav), 64)

	//RAM
	var ramTotalVal bytes.Buffer
	var ramFreeVal bytes.Buffer

	//regex
	ramreg, _ := regexp.Compile("[0-9]+")

	ramTotal.Stdout = &ramTotalVal
	ramFree.Stdout = &ramFreeVal

	ramterr := ramTotal.Run()

	if ramterr != nil {
		fmt.Println(ramterr)
	}
	ramferr := ramFree.Run()
	if ramferr != nil {
		fmt.Println(ramferr)
	}

	ramFreeF, _ := strconv.ParseFloat(ramreg.FindString(ramFreeVal.String()), 32)
	ramTotalF, _ := strconv.ParseFloat(ramreg.FindString(ramTotalVal.String()), 32)

	ramPerc := roundFloat((((ramTotalF/1000)-ramFreeF)/(ramTotalF/1000))*100, 1)

	//cpu, ram
	return Stats{cpuStat, ramPerc}
}

func getLinuxStats() Stats {
	cpu := exec.Command("top", "-b", "-n", "1")
	//CPU
	var cpuout bytes.Buffer
	cpu.Stdout = &cpuout
	cpuerr := cpu.Run()
	if cpuerr != nil {
		fmt.Println(cpuerr)
	}

	cpuavreg, _ := regexp.Compile(",[ ]*[0-9.]+ id")
	cpuavregnum, _ := regexp.Compile("[0-9.]+")
	cpuav := cpuavreg.FindString(cpuout.String())
	cpuidle, _ := strconv.ParseFloat(cpuavregnum.FindString(cpuav), 32)

	cpuused := 100 - cpuidle

	//RAM
	ramavreg, _ := regexp.Compile("MiB Mem.*?free")
	ramavfreereg, _ := regexp.Compile("[0-9.]+ free")
	ramavtotalreg, _ := regexp.Compile("[0-9.]+ total")
	ramavnumreg, _ := regexp.Compile("[0-9.]+")

	ramfreeav := ramavreg.FindString(cpuout.String())
	ramfreestr := ramavfreereg.FindString(ramfreeav)
	ramfree, _ := strconv.ParseFloat(ramavnumreg.FindString(ramfreestr), 32)

	ramtotalav := ramavreg.FindString(ramfreeav)
	ramtotalstr := ramavtotalreg.FindString(ramtotalav)
	ramtotal, _ := strconv.ParseFloat(ramavnumreg.FindString(ramtotalstr), 32)

	ramperc := roundFloat((((ramtotal)-ramfree)/(ramtotal))*100, 1)
	return Stats{cpuused, ramperc}
}

var connectedFromWebUI bool

func writeHeaders(w http.ResponseWriter, extras []string) {

	extra := strings.Join(extras, ",")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, "+extra)
}

// handlers
func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {

	global.Logger.Log("PING RECEIVED", logging.INFO)

	//fetch data of ping (nodemap)
	body, _ := ioutil.ReadAll(req.Body)

	//move nodemap to local memory
	var localnm []*global.Node
	json.Unmarshal(body, &localnm)
	//add the changed node map
	currentMasterData := global.NODE_MAP[0].Data
	global.NODE_MAP = localnm

	//only receive ping if its a follower
	if node.Rank == global.FOLLOWER {
		//print the rank of the node and wait for 2 secs
		global.Logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+strings.Join(global.GetNodeIps(), " "), logging.INFO)
		//need to check for ping here (start )
		go (*global.Node)(node).CheckForNoPingFromMaster()

	} else {
		global.Logger.Log("SOMETHINGS GONE WRONG or CONNECTED FROM WEBUI!", logging.WARNING)
		node.Rank = global.FOLLOWER
	}
	node.Pinged = time.Now()

	if len(localnm[0].Data) == 0 {
		global.NODE_MAP[0].Data = currentMasterData
		return
	}

	global.Logger.Log("UPDATED DATA ON THIS NODE!", logging.GOOD)

	global.NODE_MAP = []*global.Node{}
	for _, j := range localnm {
		global.NODE_MAP = append(global.NODE_MAP, j)
	}

}

func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, []string{"ip"})
	//find the port of client trying to connect
	ip := req.URL.Query().Get("ip")

	//if port isnt 0 (there is a valid port)
	if ip != "" {
		//add to the node map
		global.Logger.Log(ip+" has connected", logging.GOOD)
		if !global.AlreadyInNodeMap(ip) {
			global.NODE_MAP = append(global.NODE_MAP, &global.Node{Ip: ip, Pinged: time.Now(), PingCount: 0, Rank: global.FOLLOWER, Data: global.NodeData{}, Active: true})
		} else {
			//already in map
			indexOfNode := global.IndexOfNodeIpInNodeMap(ip)
			//set to active
			global.NODE_MAP[indexOfNode].Active = true
			//resend data
			global.NODE_MAP[indexOfNode].PingCount = 0
		}
	}

	//write back to the client with the rank of the node its trying to connect to
	w.Header().Add("rank", string(node.Rank))

	jsonNodeMap, _ := json.Marshal(global.NODE_MAP)
	sendData := bytes.NewBuffer(jsonNodeMap)

	w.Write(sendData.Bytes())
}

// this needs to change
func (node *Node) setDataHandler(w http.ResponseWriter, req *http.Request) {

	// This could be done using sockets rather than
	// http requests
	if node.Rank == global.MASTER {

		writeHeaders(w, []string{"all"})

		if req.Method == "POST" {

			//TODO send multiple key and values
			var dataToSet []string
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &dataToSet)

			//this will change eventually
			setFolder := dataToSet[0]
			setKey := dataToSet[1]
			setVal := dataToSet[2]
			dataType := dataToSet[3]

			value := global.JsonData{
				Key:    setKey,
				Folder: setFolder,
				Data:   setVal,
				Type:   dataType,
				Date:   time.Now(),
			}

			node.Data[setKey] = value
			global.DataChanged = true
			// add to index
			// needs to be incremental
			(*global.Node)(node).SaveDataJson()
			json.NewEncoder(w).Encode("done")
			return
		}
	}
}

func (node *Node) getAllDataHandler(w http.ResponseWriter, req *http.Request) {
	if node.Rank == global.MASTER {

		writeHeaders(w, nil)

		dataToSend := global.GetAllDataToPrint(node.Data)

		json.NewEncoder(w).Encode(dataToSend)

	}
}

func (node *Node) getDataHandler(w http.ResponseWriter, req *http.Request) {
	if node.Rank == global.MASTER {

		writeHeaders(w, []string{"key"})

		dataKey := req.URL.Query().Get("key")
		if dataKey == "" {
			dataKey = req.Header.Get("key")
		}

		getData := global.NODE_MAP[0].Data[dataKey]

		sendData := SendData{dataKey, getData}

		json.NewEncoder(w).Encode(sendData)

	}
}

func (node *Node) removeNodeHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{"ip"})

	//only do this if master
	if node.Rank == global.FOLLOWER {
		return
	}

	nodeIpToRemove := req.URL.Query().Get("ip")
	if nodeIpToRemove == "" {
		nodeIpToRemove = req.Header.Get("ip")
	}

	//this needs to be a function
	indexOfNode := global.IndexOfNodeIpInNodeMap(nodeIpToRemove)
	if indexOfNode == -1 {
		return
	}
	global.Logger.Log("removed node: "+nodeIpToRemove+" from the nodemap", logging.GOOD)

	global.NODE_MAP = append(global.NODE_MAP[:indexOfNode], global.NODE_MAP[indexOfNode+1:]...)
	json.NewEncoder(w).Encode("removed node")

}

func (node *Node) folderHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, nil)
	//get folders
	var folders []string
	for _, d := range node.Data {
		if !isInArr(folders, d.Folder) {
			folders = append(folders, d.Folder)
		}
	}
	json.NewEncoder(w).Encode(folders)
}

func (node *Node) addFolderHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{"folder"})
	//get folders
	folder := req.URL.Query().Get("folder")
	if folder == "" {
		folder = req.Header.Get("folder")
	}

	key := getHexKey()
	value := global.JsonData{
		Key:    "_firstdoc_" + key,
		Folder: folder,
		Data:   "_firstdoc",
		Type:   "text",
		Date:   time.Now(),
	}

	node.Data[key] = value
	global.DataChanged = true

	json.NewEncoder(w).Encode("done")
}

func (node *Node) getDataInFolderHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, nil)

	folder := req.URL.Query().Get("folder")

	//need to sort data by date
	//breaking change, as added new JSON field
	nodeData := make([]global.JsonData, len(node.Data))

	dataInd := 0
	for _, d := range node.Data {
		nodeData[dataInd] = d
		dataInd++
	}

	sort.Slice(nodeData, func(i, j int) bool {
		return nodeData[i].Date.Unix() > nodeData[j].Date.Unix()
	})

	var dataInFolder []SendData
	numOfItems := 0
	for i, data := range nodeData {
		key := nodeData[i].Key
		if numOfItems == 40 {
			break
		}
		if data.Folder == folder {
			sendData := SendData{key, data}
			dataInFolder = append(dataInFolder, sendData)
			numOfItems++
		}
	}
	json.NewEncoder(w).Encode(dataInFolder)

}

func (node *Node) statsHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, nil)

	os := runtime.GOOS

	osStats := map[string]func() Stats{
		"windows": getWindowsStats,
		"linux":   getLinuxStats,
	}

	stats := osStats[os]()

	json.NewEncoder(w).Encode(stats)
}

func isInArr(arr []string, needle string) bool {
	for _, s := range arr {
		if s == needle {
			return true
		}
	}
	return false
}

func nodeMapHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, []string{"all"})

	nodeMapResp := []WebNodeMap{}
	for _, n := range global.NODE_MAP {
		nodeMapResp = append(nodeMapResp, WebNodeMap{n.Ip, n.Active})
	}

	json.NewEncoder(w).Encode(nodeMapResp)
}

func (node *Node) queryHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, []string{"query"})

	queryVal := req.URL.Query().Get("query")
	if queryVal == "" {
		queryVal = req.Header.Get("query")
	}
	query.LoadIntoMemory("index/primary.json")
	queryResult, timeTaken, _, _ := query.Execute(queryVal, "")

	dataToSend := make([]SendData, 0)

	for _, res := range queryResult {

		dataJSON, _ := json.Marshal(res)
		var getJSON global.JsonData
		json.Unmarshal(dataJSON, &getJSON)

		dataToSend = append(dataToSend, SendData{getJSON.Key, getJSON})

	}

	json.NewEncoder(w).Encode(SendQueryData{dataToSend, timeTaken})
}

func SetupHandlers(node *Node) {
	//handle incoming connection (MAYBE MAKE THIS ASYNC)
	http.HandleFunc("/connect", node.connectHandler)
	http.HandleFunc("/ping", node.pingHandler)
	http.HandleFunc("/getalldata", node.getAllDataHandler)
	http.HandleFunc("/getdata", node.getDataHandler)
	http.HandleFunc("/addfolder", node.addFolderHandler)
	http.HandleFunc("/folders", node.folderHandler)
	http.HandleFunc("/removenode", node.removeNodeHandler)
	http.HandleFunc("/nodemap", nodeMapHandler)
	http.HandleFunc("/stats", node.statsHandler)

	http.HandleFunc("/getquery", node.queryHandler)

	go http.HandleFunc("/getdatainfolder", node.getDataInFolderHandler)
	go http.HandleFunc("/setdata", node.setDataHandler)
}
