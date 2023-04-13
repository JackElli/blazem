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
	// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
	// coming in, write to disk and add to the map
	WriteHeaders(w, []string{"all"})

	if req.Method != "POST" {
		JsonResponse(w, EndpointResponse{
			500,
			"Should not be getting",
			nil,
		})
		return
	}

	if node.Rank != global.MASTER {
		JsonResponse(w, EndpointResponse{
			500,
			"Should be master",
			nil,
		})
		return
	}

	var dataToAdd global.Document
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			"Cannot read request body",
			nil,
		})
		return
	}
	err = json.Unmarshal(body, &dataToAdd)
	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			"Cannot unmarshal JSON request",
			nil,
		})
		return
	}

	var document = node.transformNewDoc(dataToAdd)
	global.WriteDocToDisk(document)
	node.Data.Store(dataToAdd["key"], document)

	global.DataChanged = true
	JsonResponse(w, EndpointResponse{
		200,
		"Added document successfully",
		nil,
	})
}

func (node *Node) transformNewDoc(dataToAdd global.Document) global.Document {
	// We want to transform the document coming in, to something that is optimised and
	// info-full for retrieval
	var document = dataToAdd

	if document["type"].(string) != "text" {
		document["file_name"] = dataToAdd["file_name"]
	}

	if loadDoc, ok := node.Data.Load(dataToAdd["key"]); ok {
		return updateDocument(document, loadDoc.(global.Document))
	}
	document["date"] = time.Now().Format("2006-01-02T15:04:05")
	return document
}

func updateDocument(document global.Document, loadDoc global.Document) global.Document {
	// If we're not adding a new document, we're updating an existing one, we want
	// to keep the date the same
	document["date"] = loadDoc["date"].(string)
	return document
}
