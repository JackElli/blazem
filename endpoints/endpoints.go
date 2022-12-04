package endpoints

import (
	"bytes"
	"distributed_servers/global"
	"distributed_servers/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Node global.Node

type WebNodeMap struct {
	Ip     string
	Active bool
}
type SendData struct {
	Key  string
	Data global.JsonData
}

var connectedFromWebUI bool

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
		node.Pinged = time.Now()
		//need to check for ping here (start )
		go (*global.Node)(node).CheckForNoPingFromMaster()

	} else {
		global.Logger.Log("SOMETHINGS GONE WRONG or CONNECTED FROM WEBUI!", logging.WARNING)
		node.Rank = global.FOLLOWER
	}

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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, ip")

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
	if node.Rank == global.MASTER {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, all")

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
				Folder: setFolder,
				Data:   setVal,
				Type:   dataType,
				Date:   time.Now(),
			}

			node.Data[setKey] = value
			global.DataChanged = true

			w.Header().Set("response", "done")
			json.NewEncoder(w).Encode("done")
			return
		}
	}
}

func (node *Node) getAllDataHandler(w http.ResponseWriter, req *http.Request) {
	if node.Rank == global.MASTER {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

		dataToSend := global.GetAllDataToPrint(node.Data)

		json.NewEncoder(w).Encode(dataToSend)

	}
}

func (node *Node) getDataHandler(w http.ResponseWriter, req *http.Request) {
	if node.Rank == global.MASTER {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, key")

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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, ip")

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

func (node *Node) FolderHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	//get folders
	var folders []string
	for _, d := range node.Data {
		if !isInArr(folders, d.Folder) {
			folders = append(folders, d.Folder)
		}
	}
	json.NewEncoder(w).Encode(folders)
}

func (node *Node) getDataInFolderHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	folder := req.URL.Query().Get("folder")

	var dataInFolder []SendData
	numOfItems := 0
	for key, data := range node.Data {
		if numOfItems == 40 {
			break
		}
		if data.Folder == folder {
			sendData := SendData{key, data}
			dataInFolder = append(dataInFolder, sendData)
		}
		numOfItems++
	}
	json.NewEncoder(w).Encode(dataInFolder)

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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, all")

	nodeMapResp := []WebNodeMap{}
	for _, n := range global.NODE_MAP {
		nodeMapResp = append(nodeMapResp, WebNodeMap{n.Ip, n.Active})
	}

	json.NewEncoder(w).Encode(nodeMapResp)
}

func SetupHandlers(node *Node) {
	//handle incoming connection (MAYBE MAKE THIS ASYNC)
	http.HandleFunc("/connect", node.connectHandler)
	http.HandleFunc("/ping", node.pingHandler)
	http.HandleFunc("/getalldata", node.getAllDataHandler)
	http.HandleFunc("/getdata", node.getDataHandler)
	go http.HandleFunc("/getdatainfolder", node.getDataInFolderHandler)
	http.HandleFunc("/setdata", node.setDataHandler)
	http.HandleFunc("/folders", node.FolderHandler)
	http.HandleFunc("/removenode", node.removeNodeHandler)
	http.HandleFunc("/nodemap", nodeMapHandler)
}
