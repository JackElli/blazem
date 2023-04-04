package endpoints

import (
	"blazem/global"
	"blazem/logging"
	"blazem/query"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var connectedFromWebUI bool

func nodeMapHandler(w http.ResponseWriter, req *http.Request) {
	// Return the results of the nodemap to the client
	writeHeaders(w, []string{"all"})

	nodeMapResp := []WebNodeMap{}
	for _, n := range global.NODE_MAP {
		nodeMapResp = append(nodeMapResp, WebNodeMap{n.Ip, n.Active})
	}

	json.NewEncoder(w).Encode(nodeMapResp)
}

func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {
	// When we receive a ping from master we want to check if any data
	// has changed. If it has, we want to update our local data map. We also
	// want to check that the master is still alive with a ping from master
	// check. If we change from master to follower quickly, it's because we've
	// been added to the cluster by another node. We write all of the changed
	// data to disk
	var localTempNodes []*global.TempNode

	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &localTempNodes)
	localnm := global.UnmarshalNodeMap(localTempNodes)

	if len(localnm) == 0 {
		return
	}

	node.Pinged = time.Now()
	global.Logger.Log("PING RECEIVED", logging.INFO)

	currentMasterData := global.NODE_MAP[0].Data
	global.NODE_MAP = localnm

	if node.Rank == global.FOLLOWER {

		global.Logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+
			strings.Join(global.GetNodeIps(), " "), logging.INFO)

		go (*global.Node)(node).CheckForNoPingFromMaster()

	} else {
		global.Logger.Log("SOMETHINGS GONE WRONG or CONNECTED FROM WEBUI!",
			logging.WARNING)
		node.Rank = global.FOLLOWER
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

	global.NODE_MAP[0].Data.Range(func(key, value any) bool {
		_, err := os.Stat("data/" + key.(string))
		if os.IsNotExist(err) {
			global.WriteDocToDisk(value.(map[string]interface{}))
		}
		return true
	})
}

func (node *Node) connectHandler(w http.ResponseWriter, req *http.Request) {
	// We need to connect a node to the cluster; we check for ip, if it is already
	// in the node map, we set to active (because it must be active as it's sent a
	// connect request). If it's not in the nodemap, we add it.
	writeHeaders(w, []string{"ip"})

	ip := req.URL.Query().Get("ip")

	if ip != "" {
		if !global.AlreadyInNodeMap(ip) {
			global.NODE_MAP = append(global.NODE_MAP, &global.Node{Ip: ip, Pinged: time.Now(),
				PingCount: 0, Rank: global.FOLLOWER, Data: sync.Map{}, Active: true,
				RecentQueries: map[string]string{}, Rules: map[string]global.Rule{}})
		} else {

			indexOfNode := global.IndexOfNodeIpInNodeMap(ip)
			global.NODE_MAP[indexOfNode].Active = true
			global.NODE_MAP[indexOfNode].PingCount = 0
		}
		global.Logger.Log(ip+" has connected", logging.GOOD)
	}

	w.Header().Add("rank", string(node.Rank))

	jsonNodeMap, _ := json.Marshal(global.NODE_MAP)
	sendData := bytes.NewBuffer(jsonNodeMap)

	w.Write(sendData.Bytes())
}

func (node *Node) removeNodeHandler(w http.ResponseWriter, req *http.Request) {
	// we want to remove a node from the node map (master only). We get the index in
	// the node map of the node, then we update the node map (removing the 'to remove' node)
	// then we save the changes.
	writeHeaders(w, []string{"ip"})

	if node.Rank == global.FOLLOWER {
		return
	}

	nodeIpToRemove := req.URL.Query().Get("ip")
	if nodeIpToRemove == "" {
		nodeIpToRemove = req.Header.Get("ip")
	}

	indexOfNode := global.IndexOfNodeIpInNodeMap(nodeIpToRemove)
	if indexOfNode == -1 {
		return
	}
	global.Logger.Log("removed node: "+nodeIpToRemove+" from the nodemap", logging.GOOD)

	global.NODE_MAP = append(global.NODE_MAP[:indexOfNode], global.NODE_MAP[indexOfNode+1:]...)
	json.NewEncoder(w).Encode("removed node")

}

func (node *Node) folderHandler(w http.ResponseWriter, req *http.Request) {
	// We want to return all of the 'root folders in the data i.e every folder
	// that doesnt have a folder parent. We fetch the folder names, add them to the
	// folder map and add the corresponding document count
	writeHeaders(w, nil)
	var folders = make(map[string]Folder, 0)

	// SPLIT THESE UP INTO SEPARATE FUNCS
	node.Data.Range(func(k, value interface{}) bool {
		dataType := value.(map[string]interface{})["type"]
		if dataType == "folder" {
			var inFolder string
			var exists bool
			var backedUp bool = false
			folderKey := value.(map[string]interface{})["key"].(string)
			folderName := value.(map[string]interface{})["folderName"].(string)
			if value.(map[string]interface{})["backedUp"] != nil {
				backedUp = value.(map[string]interface{})["backedUp"].(bool)
			}
			if inFolder, exists = value.(map[string]interface{})["folder"].(string); !exists {
				inFolder = ""
			}
			folders[folderKey] = Folder{
				inFolder,
				folderKey,
				folderName,
				0,
				backedUp,
			}
		}
		return true
	})

	node.Data.Range(func(k, value interface{}) bool {
		if folder, exists := value.(map[string]interface{})["folder"].(string); exists && folder != "" {
			currDocCount := folders[folder].DocCount
			folders[folder] = Folder{
				folders[folder].Folder,
				folders[folder].Key,
				folders[folder].FolderName,
				currDocCount + 1,
				folders[folder].BackedUp,
			}
		}
		return true
	})
	for _, folder := range folders {
		if folder.Folder != "" {
			// This stores doc count
			folderData, _ := node.Data.Load(folder.Key)
			folderData.(map[string]interface{})["docCount"] = folder.DocCount
			node.Data.Store(folder.Key, folderData)
			delete(folders, folder.Key)
		}
	}
	json.NewEncoder(w).Encode(folders)
}

func (node *Node) addFolderHandler(w http.ResponseWriter, req *http.Request) {
	// We want to add a folder, we get the folder information from the user and
	// THIS NEEDS TO CHANGE
	writeHeaders(w, []string{"folder"})
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
	// We want to fetch stats based on the OS (as the method changes)
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
	// We fetch the query entered by the user, we send that to JAQL, then
	// we send back the results to the client. We also want to add these to
	// recent queries so the user can easily get back to queries they've
	// previously entered.

	writeHeaders(w, []string{"query"})

	queryVal := req.URL.Query().Get("query")
	if queryVal == "" {
		queryVal = req.Header.Get("query")
	}

	query.LoadIntoMemory(global.Node(*node))

	queryResult, timeTaken, _, _ := query.Execute(queryVal, "")
	dataToSend := make([]SendData, 0)

	for _, res := range queryResult {
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
	// Returns a list of recently entered queries
	writeHeaders(w, []string{})

	dataToSend := node.RecentQueries
	if len(dataToSend) == 0 {
		json.NewEncoder(w).Encode([]uint8{})
		return
	}
	json.NewEncoder(w).Encode(dataToSend)
}

func (node *Node) replicateFolderHandler(w http.ResponseWriter, req *http.Request) {
	// We want a way to replicate folder data to another node outside of the cluster
	// we get all of the data currently within this folder and send that
	// over HTTP to the desired node.

	writeHeaders(w, []string{})

	var replicate global.Replicate
	body, _ := ioutil.ReadAll(req.Body)
	_ = json.Unmarshal(body, &replicate)

	localFolder := replicate.LocalFolder
	ip := replicate.RemoteIp

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
