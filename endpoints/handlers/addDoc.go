package handlers

import (
	"blazem/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func AddDocHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.addDocHandler
}

func (node *Node) addDocHandler(w http.ResponseWriter, req *http.Request) {
	// This could be done using sockets rather than
	// http requests
	WriteHeaders(w, []string{"all"})

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
