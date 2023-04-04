package handlers

import (
	"blazem/global"
	"encoding/json"
	"net/http"
)

func DeleteDocHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.deleteDocHandler
}

func (node *Node) deleteDocHandler(w http.ResponseWriter, req *http.Request) {
	WriteHeaders(w, []string{})

	docKey := req.URL.Query().Get("key")
	node.Data.Delete(docKey)
	global.DataChanged = true

	json.NewEncoder(w).Encode("done")

}
