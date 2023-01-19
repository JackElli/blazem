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

type AddData struct {
	Folder string
	Key    string
	Value  string
	Format string
}

func (node *Node) getNewDoc(dataToAdd AddData) map[string]interface{} {
	// changed to map[string]interface{}
	// to add to mem index
	// and if we want custom JSON later

	// have to be really careful with casting
	// maybe have a struct like
	// key, folder, type, date and fields

	// if doc key exists
	if _, ok := node.Data.Load(dataToAdd.Key); ok {
		getDocTime, _ := node.Data.Load(dataToAdd.Key)
		docTime := getDocTime.(map[string]interface{})
		value := map[string]interface{}{
			"key":    dataToAdd.Key,
			"folder": dataToAdd.Folder,
			"value":  dataToAdd.Value,
			"type":   dataToAdd.Format,
			"date":   docTime["date"].(string),
		}
		return value
	}
	value := map[string]interface{}{
		"key":    dataToAdd.Key,
		"folder": dataToAdd.Folder,
		"value":  dataToAdd.Value,
		"type":   dataToAdd.Format,
		"date":   time.Now().Format("2006-01-02T15:04:05"),
	}
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

	//TODO send multiple key and values
	var dataToAdd AddData
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &dataToAdd)
	if err != nil {
		return
	}

	value := node.getNewDoc(dataToAdd)
	node.Data.Store(dataToAdd.Key, value)
	// node.Data[dataToAdd.Key] = value
	global.DataChanged = true

	json.NewEncoder(w).Encode("done")
	return
}

func (node *Node) deleteDocHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{})

	docKey := req.URL.Query().Get("key")
	node.Data.Delete(docKey)
	// delete(node.Data, docKey)
	global.DataChanged = true

	json.NewEncoder(w).Encode("done")

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

	// getData := global.NODE_MAP[0].Data[dataKey]
	getData, _ := global.NODE_MAP[0].Data.Load(dataKey)
	sendData := SendData{dataKey, getData.(global.JsonData)}
	json.NewEncoder(w).Encode(sendData)

}

func (node *Node) getDataInFolderHandler(w http.ResponseWriter, req *http.Request) {

	writeHeaders(w, nil)

	folder := req.URL.Query().Get("folder")

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
