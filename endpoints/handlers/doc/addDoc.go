package doc

import (
	types "blazem/domain/endpoint"
	"blazem/domain/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func NewAddDocHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return AddDocHandler
}

func AddDocHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	de := &DocEndpoint{
		Endpoint: *e,
	}
	return de.addDocHandler
}

// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
// coming in, write to disk and add to the map
func (e *DocEndpoint) addDocHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"all"})

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if req.Method != "POST" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method " + req.Method + "!=POST",
		})
		return
	}
	if e.Endpoint.Node.Rank != global.MASTER {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Should be master",
		})
		return
	}
	var dataToAdd global.Document
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Cannot read request body {" + err.Error() + "}",
		})
		return
	}
	err = json.Unmarshal(body, &dataToAdd)
	if err != nil {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Cannot unmarshal JSON request {" + err.Error() + "}",
		})
		return
	}

	document := TransformNewDoc(e.Endpoint.Node, dataToAdd)
	global.WriteDocToDisk(document)
	e.Endpoint.Node.Data.Store(dataToAdd["key"], document)
	global.DataChanged = true

	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Added document successfully",
	})
}

// We want to transform the document coming in, to something that is optimised and
// info-full for retrieval
func TransformNewDoc(node *global.Node, dataToAdd global.Document) global.Document {
	document := dataToAdd
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
