package handlers

import (
	"blazem/global"
	"net/http"
)

func DeleteDocHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.deleteDocHandler
}

// We want to delete a document from Blazem
func (node *Node) deleteDocHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, []string{})

	if req.Method != "DELETE" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}

	var docKey = req.URL.Query().Get("key")
	if docKey == "" {
		JsonResponse(w, EndpointResponse{
			500,
			"Doc key not provided",
			nil,
		})
		return
	}

	var _, docFound = node.Data.Load(docKey)
	if !docFound {
		JsonResponse(w, EndpointResponse{
			404,
			"Cannot delete document, as it's not found",
			nil,
		})
		return
	}
	node.Data.Delete(docKey)
	global.DataChanged = true
	JsonResponse(w, EndpointResponse{
		200,
		"Successfully deleted doc",
		nil,
	})
}
