package endpoints

import (
	"distributed_servers/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

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
			// if key is not in index
			// we can append
			if !global.IsDocInIndex(setKey) {
				(*global.Node)(node).AppendDataJson(setKey, value)
			} else {
				// if key is ALREADY in index
				// we need to replace it
				(*global.Node)(node).ReplaceDataJson(setKey, value)
			}

			// (*global.Node)(node).SaveDataJson()
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
	//only do this if master
	if node.Rank == global.FOLLOWER {
		return
	}

	writeHeaders(w, []string{"key"})

	dataKey := req.URL.Query().Get("key")
	if dataKey == "" {
		dataKey = req.Header.Get("key")
	}

	getData := global.NODE_MAP[0].Data[dataKey]
	sendData := SendData{dataKey, getData}
	json.NewEncoder(w).Encode(sendData)

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
