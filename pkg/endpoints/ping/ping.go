package ping

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	blazem_node "blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type PingMgr struct {
	Router    *mux.Router
	Node      *node.Node
	Responder responder.Responder
}

func NewPingMgr(router *mux.Router, node *node.Node, responder responder.Responder) *PingMgr {
	return &PingMgr{
		Router:    router,
		Node:      node,
		Responder: responder,
	}
}

// When we receive a ping from master we want to check if any data
// has changed. If it has, we want to update our local data map. We also
// want to check that the master is still alive with a ping from master
// check. If we change from master to follower quickly, it's because we've
// been added to the cluster by another node. We write all of the changed
// data to disk
func (e *PingMgr) Ping() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		localTempNodes := make([]*blazem_node.TempNode, 0)
		json.NewDecoder(req.Body).Decode(&localTempNodes)

		localnm := node.UnmarshalNodeMap(localTempNodes)
		if len(localnm) == 0 {
			e.Responder.Error(w, 500, errors.New("No nodes found to marshal"))
			return
		}
		e.Node.Pinged = time.Now()
		logger.Logger.Info("Ping received.")
		currentMasterData := e.Node.NodeMap[0].Data
		e.Node.NodeMap = localnm

		if e.Node.Rank == global.FOLLOWER {
			go e.Node.CheckForNoPingFromMaster()
		} else {
			e.Node.Rank = global.FOLLOWER
		}

		if types.LenOfSyncMap(localnm[0].Data) == 0 {
			e.Node.NodeMap[0].Data = currentMasterData
			e.Responder.Respond(w, 200, types.EndpointResponse{
				Msg: "Successful ping",
			})
			return
		}
		e.UpdateData(localnm)
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg: "Successful ping",
		})
	}
}

// We want to add any nodes not in the nodemap to the nodemap.
// We also want to write doc to disk if it doesn't exist.
func (e *PingMgr) UpdateData(localnm []*node.Node) {
	e.Node.NodeMap = make([]*blazem_node.Node, 0)
	for _, j := range localnm {
		e.Node.NodeMap = append(e.Node.NodeMap, j)
	}
	e.Node.NodeMap[0].Data.Range(func(key, value any) bool {
		_, err := os.Stat("data/" + key.(string))
		if os.IsNotExist(err) {
			e.Node.WriteDocToDisk(key.(string), value)
		}
		return true
	})
}

func (e *PingMgr) Register() {
	e.Router.HandleFunc("/ping", e.Ping()).Methods("POST")
}
