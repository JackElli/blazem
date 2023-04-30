package ping

import (
	types "blazem/domain/endpoint"
	global_types "blazem/domain/global"
	"blazem/global"
	"blazem/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func NewPingHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return PingHandler
}

func PingHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &PingEndpoint{
		Endpoint: *e,
	}
	return pe.pingHandler
}

// When we receive a ping from master we want to check if any data
// has changed. If it has, we want to update our local data map. We also
// want to check that the master is still alive with a ping from master
// check. If we change from master to follower quickly, it's because we've
// been added to the cluster by another node. We write all of the changed
// data to disk
func (e *PingEndpoint) pingHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Wrong method",
		})
		return
	}
	localTempNodes := make([]*global.TempNode, 0)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Error reading request body {" + err.Error() + "}",
		})
		return
	}
	err = json.Unmarshal(body, &localTempNodes)
	if err != nil {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "Error unmarshalling request body {" + err.Error() + "}",
		})
		return
	}
	localnm := global.UnmarshalNodeMap(localTempNodes)
	if len(localnm) == 0 {
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 500,
			Msg:  "No nodes found to marshal",
		})
		return
	}
	e.Endpoint.Node.Pinged = time.Now()
	global.Logger.Log("PING RECEIVED", logging.INFO)
	currentMasterData := global.NODE_MAP[0].Data
	global.NODE_MAP = localnm

	if e.Endpoint.Node.Rank == global.FOLLOWER {
		go e.Endpoint.Node.CheckForNoPingFromMaster()
	} else {
		e.Endpoint.Node.Rank = global.FOLLOWER
	}

	if global_types.LenOfSyncMap(localnm[0].Data) == 0 {
		global.NODE_MAP[0].Data = currentMasterData
		e.Endpoint.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successful ping",
		})
		return
	}
	UpdateData(e.Endpoint.Node, localnm)
	e.Endpoint.Respond(w, types.EndpointResponse{
		Code: 200,
		Msg:  "Successful ping",
	})
}

// We want to add any nodes not in the nodemap to the nodemap.
// We also want to write doc to disk if it doesn't exist.
func UpdateData(node *global.Node, localnm []*global.Node) {
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
