package handlers

import (
	"blazem/global"
	"blazem/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func PingHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.pingHandler
}

func (node *Node) pingHandler(w http.ResponseWriter, req *http.Request) {
	// When we receive a ping from master we want to check if any data
	// has changed. If it has, we want to update our local data map. We also
	// want to check that the master is still alive with a ping from master
	// check. If we change from master to follower quickly, it's because we've
	// been added to the cluster by another node. We write all of the changed
	// data to disk
	var localTempNodes []*global.TempNode

	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &localTempNodes)
	var localnm = global.UnmarshalNodeMap(localTempNodes)

	if len(localnm) == 0 {
		return
	}

	node.Pinged = time.Now()
	global.Logger.Log("PING RECEIVED", logging.INFO)

	var currentMasterData = global.NODE_MAP[0].Data
	global.NODE_MAP = localnm

	if node.Rank == global.FOLLOWER {
		global.Logger.Log(string(node.Rank)+" at "+node.Ip+" nodemap: "+
			strings.Join(global.GetNodeIps(), " "), logging.INFO)

		go (*global.Node)(node).CheckForNoPingFromMaster()
	} else {
		global.Logger.Log("SOMETHINGS GONE WRONG or CONNECTED FROM WEBUI!",
			logging.WARNING)
		node.Rank = global.FOLLOWER
	}

	if lenOfSyncMap(localnm[0].Data) == 0 {
		global.NODE_MAP[0].Data = currentMasterData
		return
	}

	node.updateData(localnm)
}

func (node *Node) updateData(localnm []*global.Node) {
	global.Logger.Log("UPDATED DATA ON THIS NODE!", logging.GOOD)

	global.NODE_MAP = []*global.Node{}
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
