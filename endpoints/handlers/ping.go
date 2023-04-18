package handlers

import (
	"blazem/global"
	"blazem/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func PingHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.pingHandler
}

// When we receive a ping from master we want to check if any data
// has changed. If it has, we want to update our local data map. We also
// want to check that the master is still alive with a ping from master
// check. If we change from master to follower quickly, it's because we've
// been added to the cluster by another node. We write all of the changed
// data to disk
func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		JsonResponse(w, EndpointResponse{
			500,
			"Wrong method",
			nil,
		})
		return
	}
	var localTempNodes = make([]*global.TempNode, 0)

	var body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			"Error reading request body {" + err.Error() + "}",
			nil,
		})
		return
	}
	err = json.Unmarshal(body, &localTempNodes)
	if err != nil {
		JsonResponse(w, EndpointResponse{
			500,
			"Error unmarshalling request body {" + err.Error() + "}",
			nil,
		})
		return
	}
	var localnm = global.UnmarshalNodeMap(localTempNodes)
	if len(localnm) == 0 {
		JsonResponse(w, EndpointResponse{
			500,
			"No nodes found to marshal",
			nil,
		})
		return
	}

	node.Pinged = time.Now()
	global.Logger.Log("PING RECEIVED", logging.INFO)

	var currentMasterData = global.NODE_MAP[0].Data
	global.NODE_MAP = localnm

	if node.Rank == global.FOLLOWER {
		go (*global.Node)(node).CheckForNoPingFromMaster()
	} else {
		node.Rank = global.FOLLOWER
	}

	if lenOfSyncMap(localnm[0].Data) == 0 {
		global.NODE_MAP[0].Data = currentMasterData
		JsonResponse(w, EndpointResponse{
			200,
			"Successfull ping",
			nil,
		})
		return
	}
	node.updateData(localnm)
	JsonResponse(w, EndpointResponse{
		200,
		"Successfull ping",
		nil,
	})
}

// We want to add any nodes not in the nodemap to the nodemap.
// We also want to write doc to disk if it doesn't exist.
func (node *Node) updateData(localnm []*global.Node) {
	global.NODE_MAP = make([]*global.Node, 0)
	for _, j := range localnm {
		global.NODE_MAP = append(global.NODE_MAP, j)
	}

	global.NODE_MAP[0].Data.Range(func(key, value any) bool {
		_, err := os.Stat("data/" + key.(string))
		if os.IsNotExist(err) {
			global.WriteDocToDisk(value.(map[string]interface{}))
		}
		return true
	})
}
