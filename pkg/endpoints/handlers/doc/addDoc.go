package doc

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
// coming in, write to disk and add to the map
func AddDoc(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank != global.MASTER {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Should be master",
			})
			return
		}
		var dataToAdd global.Document
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Cannot read request body {" + err.Error() + "}",
			})
			return
		}
		err = json.Unmarshal(body, &dataToAdd)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Cannot unmarshal JSON request {" + err.Error() + "}",
			})
			return
		}

		document := TransformNewDoc(e.Node, dataToAdd)
		global.WriteDocToDisk(document)
		e.Node.Data.Store(dataToAdd["key"], document)
		global.DataChanged = true

		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Added document successfully",
		})
	}
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
