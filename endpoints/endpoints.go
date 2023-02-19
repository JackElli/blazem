package endpoints

import (
	"blazem/global"
	"blazem/logging"
	"blazem/query"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Node global.Node

type WebNodeMap struct {
	Ip     string `json:"ip"`
	Active bool   `json:"active"`
}

type SendData struct {
	Key  string                 `json:"key"`
	Data map[string]interface{} `json:"data"`
}

type Stats struct {
	Cpu float64 `json:"cpu"`
	Ram float64 `json:"ram"`
}

type SendQueryData struct {
	Docs      []SendData `json:"docs"`
	TimeTaken int64      `json:"timeTaken"`
}

type Folder struct {
	FolderName string `json:"folderName"`
	DocCount   int    `json:"docCount"`
}

var connectedFromWebUI bool

func writeHeaders(w http.ResponseWriter, extras []string) {
	extra := strings.Join(extras, ",")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, "+extra)
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

func lenOfSyncMap(mp sync.Map) int {
	var i int
	mp.Range(func(key any, value any) bool {
		i++
		return true
	})
	return i
}

func isInArr(arr []string, needle string) bool {
	for _, s := range arr {
		if s == needle {
			return true
		}
	}
	return false
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

func getMacStats() Stats {
	return Stats{
		1.1,
		2.3,
	}
}

func nodeMapHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, []string{"all"})

	nodeMapResp := []WebNodeMap{}
	for _, n := range global.NODE_MAP {
		nodeMapResp = append(nodeMapResp, WebNodeMap{n.Ip, n.Active})
	}

	json.NewEncoder(w).Encode(nodeMapResp)
}

func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {

	global.Logger.Log("PING RECEIVED", logging.INFO)

	//fetch data of ping (nodemap)
	body, _ := ioutil.ReadAll(req.Body)

	//move nodemap to local memory
	var localTempNodes []*global.TempNode
	json.Unmarshal(body, &localTempNodes)

	// this can be functionizated
	localnm := global.UnmarshalNodeMap(localTempNodes)

	//add the changed node map
	currentMasterData := global.NODE_MAP[0].Data
	global.NODE_MAP = localnm

	//only receive ping if its a follower
	if node.Rank == global.FOLLOWER {
		//print the rank of the node and wait for 2 secs
		global.Logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+
			strings.Join(global.GetNodeIps(), " "), logging.INFO)
		//need to check for ping here (start )
		go (*global.Node)(node).CheckForNoPingFromMaster()

	} else {
		global.Logger.Log("SOMETHINGS GONE WRONG or CONNECTED FROM WEBUI!", logging.WARNING)
		node.Rank = global.FOLLOWER
	}
	node.Pinged = time.Now()

	if len(localnm) == 0 {
		return
	}

	if lenOfSyncMap(localnm[0].Data) == 0 {
		global.NODE_MAP[0].Data = currentMasterData
		return
	}

	global.Logger.Log("UPDATED DATA ON THIS NODE!", logging.GOOD)

	global.NODE_MAP = []*global.Node{}
	for _, j := range localnm {
		global.NODE_MAP = append(global.NODE_MAP, j)
	}

	// write all changed docs to disk
	global.NODE_MAP[0].Data.Range(func(key, value any) bool {
		_, err := os.Stat("data/" + key.(string))
		if os.IsNotExist(err) {
			global.WriteDocToDisk(value.(map[string]interface{}))
		}
		return true
	})
}

func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, []string{"ip"})
	//find the port of client trying to connect
	ip := req.URL.Query().Get("ip")

	//if port isnt 0 (there is a valid port)
	if ip != "" {
		// add to the node map
		// a fresh node
		global.Logger.Log(ip+" has connected", logging.GOOD)
		if !global.AlreadyInNodeMap(ip) {
			global.NODE_MAP = append(global.NODE_MAP, &global.Node{Ip: ip, Pinged: time.Now(),
				PingCount: 0, Rank: global.FOLLOWER, Data: sync.Map{}, Active: true,
				RecentQueries: map[string]string{}, Rules: map[string]global.Rule{}})
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

	// this can be improved
	var folderNames = make([]string, 0)
	var folderMap = make(map[string]Folder, 0)
	var folders = make([]Folder, 0)

	// get all folder names
	node.Data.Range(func(k, value interface{}) bool {
		if !isInArr(folderNames, value.(map[string]interface{})["folder"].(string)) {
			folderNames = append(folderNames, value.(map[string]interface{})["folder"].(string))
		}
		return true
	})

	// set up folders
	for _, folder := range folderNames {
		folderMap[folder] = Folder{
			folder,
			0,
		}
	}

	// get doc count within folder
	// (used a map for easy access)
	node.Data.Range(func(k, value interface{}) bool {
		folder := value.(map[string]interface{})["folder"].(string)
		currDocCount := folderMap[folder].DocCount
		folderMap[folder] = Folder{
			folder,
			currDocCount + 1,
		}
		return true
	})

	// push to client as list of obj

	for _, folder := range folderMap {
		folders = append(folders, folder)
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
	value := map[string]interface{}{
		"key":    "_firstdoc_" + key,
		"folder": folder,
		"value":  "_firstdoc",
		"type":   "text",
		"date":   time.Now().Format("2006-01-02T15:04:05"),
	}

	node.Data.Store(key, value)
	global.DataChanged = true

	json.NewEncoder(w).Encode("done")
}

func (node *Node) statsHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, nil)

	os := runtime.GOOS

	osStats := map[string]func() Stats{
		"windows": getWindowsStats,
		"linux":   getLinuxStats,
		"darwin":  getMacStats,
	}

	stats := osStats[os]()

	json.NewEncoder(w).Encode(stats)
}

func (node *Node) queryHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, []string{"query"})

	queryVal := req.URL.Query().Get("query")
	if queryVal == "" {
		queryVal = req.Header.Get("query")
	}

	// there is some optimisation we can do
	// here
	query.LoadIntoMemory(global.Node(*node))

	// TODO error handling
	queryResult, timeTaken, _, _ := query.Execute(queryVal, "")
	dataToSend := make([]SendData, 0)

	for _, res := range queryResult {
		// we want to not send lots of data
		// let the client fetch that
		// whem rendered
		if res["type"] != "text" {
			res["value"] = "file"
		}
		dataJSON, _ := json.Marshal(res)
		var getJSON global.JsonData
		json.Unmarshal(dataJSON, &getJSON)
		dataToSend = append(dataToSend, SendData{getJSON["key"].(string), getJSON})
	}

	node.RecentQueries[queryVal] = time.Now().Format("2006-01-02 15:04:05")
	json.NewEncoder(w).Encode(SendQueryData{dataToSend, timeTaken})
}

func (node *Node) getRecentQueriesHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{})

	dataToSend := node.RecentQueries
	if len(dataToSend) == 0 {
		json.NewEncoder(w).Encode([]uint8{})
		return
	}
	json.NewEncoder(w).Encode(dataToSend)
}

func (node *Node) replicateFolderHandler(w http.ResponseWriter, req *http.Request) {

	// THIS IS UNSAFE!!!
	// NEEDS PERMS
	writeHeaders(w, []string{})

	var replicate global.Replicate
	body, _ := ioutil.ReadAll(req.Body)
	_ = json.Unmarshal(body, &replicate)

	localFolder := replicate.LocalFolder
	ip := replicate.RemoteIp

	// THIS IS QUITE INEFFICIENT
	node.Data.Range(func(key, value any) bool {
		doc := value.(map[string]interface{})
		if doc["folder"].(string) == localFolder {
			dataToJSON, _ := json.Marshal(value.(map[string]interface{}))
			dataToSend := bytes.NewReader(dataToJSON)
			http.Post("http://"+ip+"/addDoc", "application/json", dataToSend)
		}
		return true
	})

	json.NewEncoder(w).Encode("done")
	return
}

func SetupHandlers(node *Node) {

	var handlers = map[string]map[string]func(http.ResponseWriter, *http.Request){
		"sync": {
			"connect":          node.connectHandler,
			"getDoc":           node.getDocHandler,
			"addFolder":        node.addFolderHandler,
			"folders":          node.folderHandler,
			"removeNode":       node.removeNodeHandler,
			"stats":            node.statsHandler,
			"nodemap":          nodeMapHandler,
			"getQuery":         node.queryHandler,
			"getRecentQueries": node.getRecentQueriesHandler,
			"addRule":          node.addRuleHandler,
			"removeRule":       node.removeRuleHandler,
			"runRule":          node.runRuleHandler,
			"getRules":         node.getRulesHandler,
		},
		"async": {
			"getDataInFolder": node.getDataInFolderHandler,
			"addDoc":          node.addDocHandler,
			"deleteDoc":       node.deleteDocHandler,
			"ping":            node.pingHandler,
			"replicateFolder": node.replicateFolderHandler,
		},
	}

	//sync
	for fncType, handlerMap := range handlers {
		for end, fnc := range handlerMap {
			if fncType == "sync" {
				http.HandleFunc("/"+end, fnc)
				continue
			}
			go http.HandleFunc("/"+end, fnc)
		}
	}
}
