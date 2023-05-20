package doc

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/node"
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
		var dataToAdd map[string]interface{}
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
		e.Node.WriteDocToDisk(document)
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
func TransformNewDoc(node *node.Node,
	dataToAdd map[string]interface{}) map[string]interface{} {
	document := dataToAdd

	document["date"] = time.Now().Format("2006-01-02T15:04:05")

	_, typeExist := document["type"]
	if !typeExist {
		// should read file type
		document["type"] = "text"
	}

	loadDoc, ok := node.Data.Load(dataToAdd["key"])
	if ok {
		return updateDocument(document, loadDoc.(map[string]interface{}))
	}

	return document
}

// If we're not adding a new document, we're updating an existing one, we want
// to keep the date the same
func updateDocument(document map[string]interface{},
	loadDoc map[string]interface{}) map[string]interface{} {
	document["date"] = loadDoc["date"].(string)
	return document
}
