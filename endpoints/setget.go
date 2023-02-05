package endpoints

import (
	"blazem/global"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

func (node *Node) getNewDoc(dataToAdd map[string]interface{}) map[string]interface{} {
	// changed to map[string]interface{}
	// to add to mem index
	// and if we want custom JSON later
	value := dataToAdd

	// if its not text, set the value
	// to not text so we dont add
	// the image/file to memory
	if value["type"].(string) != "text" {
		// value["value"] = "file"
		// set the file_name field if its not
		// just text
		value["file_name"] = dataToAdd["file_name"]
	}
	// have to be really careful with casting
	// maybe have a struct like
	// key, folder, type, date and fields
	// if doc key exists dont update date
	if _, ok := node.Data.Load(dataToAdd["key"]); ok {
		getDocTime, _ := node.Data.Load(dataToAdd["key"])
		docTime := getDocTime.(map[string]interface{})
		value["date"] = docTime["date"].(string)
		return value
	}
	value["date"] = time.Now().Format("2006-01-02T15:04:05")
	return value
}

func (node *Node) addDocHandler(w http.ResponseWriter, req *http.Request) {
	// This could be done using sockets rather than
	// http requests
	writeHeaders(w, []string{"all"})

	if node.Rank != global.MASTER {
		return
	}

	if req.Method != "POST" {
		return
	}

	// changed to map[string]interface{}
	// to allow any fields to be passed
	var dataToAdd map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &dataToAdd)
	if err != nil {
		return
	}

	// a bit of processing on the incoming
	// doc
	value := node.getNewDoc(dataToAdd)
	// if the data is not just text
	// we need to add it to disk
	// this is for backup
	global.WriteDocToDisk(value)

	// so we can replicate changes
	global.DataChanged = true

	// store the doc
	node.Data.Store(dataToAdd["key"], value)

	json.NewEncoder(w).Encode("done")
	return
}

func (node *Node) deleteDocHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{})

	docKey := req.URL.Query().Get("key")
	node.Data.Delete(docKey)
	global.DataChanged = true

	json.NewEncoder(w).Encode("done")

}

func (node *Node) getDocHandler(w http.ResponseWriter, req *http.Request) {
	//only do this if master
	if node.Rank == global.FOLLOWER {
		return
	}

	writeHeaders(w, []string{"key"})

	dataKey := req.URL.Query().Get("key")

	if dataKey == "" {
		dataKey = req.Header.Get("key")
	}
	getData, _ := global.NODE_MAP[0].Data.Load(dataKey)

	// make a copy, so if we need to overwrite value
	// if doesnt overwrite in the node
	sendDataJson := getData.(map[string]interface{})

	// do we need to load from disk?
	if getData.(map[string]interface{})["type"] != "text" {
		data, _ := ioutil.ReadFile("data/" + dataKey)
		var dataJSON global.JsonData
		json.Unmarshal(data, &dataJSON)
		sendDataJson["value"] = dataJSON["value"].(string)
	}
	sendData := SendData{dataKey, sendDataJson}
	json.NewEncoder(w).Encode(sendData)

}

func (node *Node) getDataInFolderHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, nil)

	folder := req.URL.Query().Get("folder")
	user := req.URL.Query().Get("user")

	if user != "jack" {
		json.NewEncoder(w).Encode("no auth")
		return
	}

	// need to sort data by date
	// breaking change, as added new JSON field
	nodeData := make([]map[string]interface{}, lenOfSyncMap(node.Data))

	dataInd := 0
	node.Data.Range(func(key, value interface{}) bool {
		nodeData[dataInd] = value.(map[string]interface{})
		dataInd++
		return true
	})

	// not sure why this is like this?
	sort.Slice(nodeData, func(i, j int) bool {
		if _, convOk := nodeData[i]["date"].(time.Time); !convOk {
			dateI, errI := time.Parse("2006-01-02T15:04:05", nodeData[i]["date"].(string))
			dateJ, errJ := time.Parse("2006-01-02T15:04:05", nodeData[j]["date"].(string))
			if errI != nil || errJ != nil {
				log.Fatal(errI)
			}
			return dateI.Unix() > dateJ.Unix()
		}
		return nodeData[i]["date"].(time.Time).Unix() > nodeData[j]["date"].(time.Time).Unix()
	})

	var dataInFolder []SendData
	numOfItems := 0
	for i, data := range nodeData {
		key := nodeData[i]["key"].(string)

		// push the word file instead
		// of file data so we save
		// on loading times
		if data["type"] != "text" {
			data["value"] = "file"
		}

		if numOfItems == 40 {
			break
		}

		if data["folder"] == folder {
			sendData := SendData{key, data}
			dataInFolder = append(dataInFolder, sendData)
			numOfItems++
		}
	}

	json.NewEncoder(w).Encode(dataInFolder)
}
