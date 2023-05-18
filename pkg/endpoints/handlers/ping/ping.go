package ping

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	blazem_node "blazem/pkg/domain/node"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// When we receive a ping from master we want to check if any data
// has changed. If it has, we want to update our local data map. We also
// want to check that the master is still alive with a ping from master
// check. If we change from master to follower quickly, it's because we've
// been added to the cluster by another node. We write all of the changed
// data to disk
func Ping(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		localTempNodes := make([]*blazem_node.TempNode, 0)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Error reading request body {" + err.Error() + "}",
			})
			return
		}
		err = json.Unmarshal(body, &localTempNodes)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Error unmarshalling request body {" + err.Error() + "}",
			})
			return
		}
		localnm := node.UnmarshalNodeMap(localTempNodes)
		if len(localnm) == 0 {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No nodes found to marshal",
			})
			return
		}
		e.Node.Pinged = time.Now()
		logger.Logger.Info("PING RECEIVED")
		currentMasterData := e.Node.NodeMap[0].Data
		e.Node.NodeMap = localnm

		if e.Node.Rank == global.FOLLOWER {
			go e.Node.CheckForNoPingFromMaster()
		} else {
			e.Node.Rank = global.FOLLOWER
		}

		if types.LenOfSyncMap(localnm[0].Data) == 0 {
			e.Node.NodeMap[0].Data = currentMasterData
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 200,
				Msg:  "Successful ping",
			})
			return
		}
		UpdateData(e.Node, localnm)
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successful ping",
		})
	}
}

// We want to add any nodes not in the nodemap to the nodemap.
// We also want to write doc to disk if it doesn't exist.
func UpdateData(node *node.Node, localnm []*node.Node) {
	node.NodeMap = make([]*blazem_node.Node, 0)
	for _, j := range localnm {
		node.NodeMap = append(node.NodeMap, j)
	}
	node.NodeMap[0].Data.Range(func(key, value any) bool {
		_, err := os.Stat("data/" + key.(string))
		if os.IsNotExist(err) {
			node.WriteDocToDisk(value.(map[string]interface{}))
		}
		return true
	})
}
