package doc

import (
	types "blazem/domain/endpoint"
	"blazem/domain/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func NewGetDocHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return GetDocHandler
}

func GetDocHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	de := &DocEndpoint{
		Endpoint: *e,
	}
	return de.getDocHandler
}

// We want to fetch a document and return it to the user
func (e *DocEndpoint) getDocHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"key"})

	if req.Method != "GET" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}

	if e.Endpoint.Node.Rank == global.FOLLOWER {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Cannot fetch doc from a follower node",
		})
		return
	}
	var dataKey = req.URL.Query().Get("key")
	if dataKey == "" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Doc key not provided",
		})
		return
	}
	var getData, ok = global.NODE_MAP[0].Data.Load(dataKey)
	if !ok {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 404,
			Msg:  "Doc not found",
		})
		return
	}
	sendDataJson := formatData(getData.(global.Document), dataKey)
	sendData := types.SendData{
		Key:  dataKey,
		Data: sendDataJson,
	}
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successfully retrieved doc",
		Data: sendData,
	})
}

// We want to do a bit of manipulation to the document for instance
// we want to read from disk if it's not a text file or a folder
func formatData(getData global.Document, dataKey string) global.Document {
	if getData["type"] == "text" {
		return getData
	}

	if getData["type"] == "folder" {
		return getData
	}

	data, _ := ioutil.ReadFile("data/" + dataKey)
	var dataJSON global.JsonData
	json.Unmarshal(data, &dataJSON)
	getData["value"] = dataJSON["value"].(string)
	return getData
}
