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
	WriteHeaders(w, []string{"key"})

	if req.Method != "GET" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	if node.Rank == global.FOLLOWER {
		JsonResponse(w, EndpointResponse{
			500,
			"Cannot fetch doc from a follower node",
			nil,
		})
		return
	}
	var dataKey = req.URL.Query().Get("key")
	if dataKey == "" {
		JsonResponse(w, EndpointResponse{
			500,
			"Doc key not provided",
			nil,
		})
		return
	}
	var getData, ok = global.NODE_MAP[0].Data.Load(dataKey)
	if !ok {
		JsonResponse(w, EndpointResponse{
			404,
			"Doc not found",
			nil,
		})
		return
	}
	var sendDataJson = formatData(getData.(global.Document), dataKey)
	var sendData = SendData{dataKey, sendDataJson}

	JsonResponse(w, EndpointResponse{
		200,
		"Successfully retrieved doc",
		sendData,
	})
}

func formatData(getData global.Document, dataKey string) global.Document {
	if getData["type"] == "text" {
		return getData
	}

	if getData["type"] == "folder" {
		return getData
	}

	var data, _ = ioutil.ReadFile("data/" + dataKey)
	var dataJSON global.JsonData
	json.Unmarshal(data, &dataJSON)
	getData["value"] = dataJSON["value"].(string)
	return getData
}
