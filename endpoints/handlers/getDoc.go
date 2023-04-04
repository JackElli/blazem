package handlers

import (
	"blazem/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetDocHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.getDocHandler
}
func (node *Node) getDocHandler(w http.ResponseWriter, req *http.Request) {
	// We want to fetch a document and return it to the user
	if node.Rank == global.FOLLOWER {
		return
	}
	WriteHeaders(w, []string{"key"})
	dataKey := req.URL.Query().Get("key")
	if dataKey == "" {
		dataKey = req.Header.Get("key")
	}
	getData, _ := global.NODE_MAP[0].Data.Load(dataKey)
	sendDataJson := getData.(map[string]interface{})
	if getData.(map[string]interface{})["type"] != "text" {
		data, _ := ioutil.ReadFile("data/" + dataKey)
		var dataJSON global.JsonData
		json.Unmarshal(data, &dataJSON)
		sendDataJson["value"] = dataJSON["value"].(string)
	}
	sendData := SendData{dataKey, sendDataJson}
	json.NewEncoder(w).Encode(sendData)
}
