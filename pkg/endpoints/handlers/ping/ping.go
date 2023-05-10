package ping

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/logging"
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
func Ping(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Wrong method",
			})
			return
		}
		localTempNodes := make([]*global.TempNode, 0)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Error reading request body {" + err.Error() + "}",
			})
			return
		}
		err = json.Unmarshal(body, &localTempNodes)
		if err != nil {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Error unmarshalling request body {" + err.Error() + "}",
			})
			return
		}
		localnm := global.UnmarshalNodeMap(localTempNodes)
		if len(localnm) == 0 {
			r.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No nodes found to marshal",
			})
			return
		}
		r.Node.Pinged = time.Now()
		global.Logger.Log("PING RECEIVED", logging.INFO)
		currentMasterData := global.NODE_MAP[0].Data
		global.NODE_MAP = localnm

		if r.Node.Rank == global.FOLLOWER {
			go r.Node.CheckForNoPingFromMaster()
		} else {
			r.Node.Rank = global.FOLLOWER
		}

		if types.LenOfSyncMap(localnm[0].Data) == 0 {
			global.NODE_MAP[0].Data = currentMasterData
			r.Respond(w, types.EndpointResponse{
				Code: 200,
				Msg:  "Successful ping",
			})
			return
		}
		UpdateData(r.Node, localnm)
		r.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Successful ping",
		})
	}
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
