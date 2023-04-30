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

// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
// coming in, write to disk and add to the map
func (node *Node) addDocHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, []string{"all"})

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if req.Method != "POST" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method " + req.Method + "!=POST",
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
			"Cannot read request body {" + err.Error() + "}",
			nil,
		})
		return
	}
	err = json.Unmarshal(body, &dataToAdd)
	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			"Cannot unmarshal JSON request {" + err.Error() + "}",
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

// We want to transform the document coming in, to something that is optimised and
// info-full for retrieval
func (node *Node) transformNewDoc(dataToAdd global.Document) global.Document {
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

// If we're not adding a new document, we're updating an existing one, we want
// to keep the date the same
func updateDocument(document global.Document, loadDoc global.Document) global.Document {
	document["date"] = loadDoc["date"].(string)
	return document
}
