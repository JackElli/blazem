package node

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/logger"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type INode interface {
	PingRetry(sendData *bytes.Buffer) bool
	PingEachConnection(jsonNodeMap []byte)
	Ping()
	CheckForNoPingFromMaster()
	ReadFromLocal()
	IsNextInLine() bool
	CheckIfDataChanged() []byte
	WriteDocToDisk(value map[string]interface{})
	GetNodeIps() []string
	AlreadyInNodeMap(ip string) bool
	IndexOfNodeIpInNodeMap(ip string) int
	GetNodeMapWithoutData() []*Node
	MarshalNodeMap() []*TempNode
	SetNodeMasterAttrs()
	PickPort(ip string) error
	tryListen(ip string)
	IndexOfNodeInNodeMap() int
	GetNodeDatas() []sync.Map
	GetLocalIp() string
	SetupLogger() error
}

type Node struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          global.Rank
	Data          sync.Map
	Active        bool
	RecentQueries map[string]string //time
	Rules         map[string]global.Rule
	NodeMap       []*Node
}

type TempNode struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          global.Rank
	Data          map[string]interface{}
	Active        bool
	RecentQueries map[string]string //time
	Rules         map[string]global.Rule
}

func NewNode() *Node {
	return &Node{
		Ip:            "",
		Pinged:        time.Now(),
		PingCount:     0,
		Rank:          global.FOLLOWER,
		Data:          sync.Map{},
		Active:        true,
		RecentQueries: map[string]string{},
		Rules:         map[string]global.Rule{},
	}
}

// PingRetry retries the ping 3 times and if afer 3 pings there's no response,
// node is 'paused'
func (n *Node) PingRetry(sendData *bytes.Buffer) bool {
	if n == nil || n.Active == false {
		return false
	}
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		_, err := http.Post("http://"+n.Ip+"/ping", "application/json", sendData)
		if err == nil {
			return true
		}
		logger.Logger.Info("Pinging again.")
	}
	return false
}

// We want to ping each follower node to make sure they know, the
// master is still active
func (node *Node) PingEachConnection(jsonNodeMap []byte) {
	for _, n := range node.NodeMap {
		go func(loopn *Node) {
			if loopn.Ip == node.Ip {
				return
			}
			if loopn.Active == false {
				if time.Now().Second()%7 != 0 && time.Now().Second()%8 != 0 {
					return
				}
			}
			sendData := bytes.NewBuffer(jsonNodeMap)
			logger.Logger.Info("Pinging " + loopn.Ip)

			_, err := http.Post("http://"+loopn.Ip+"/ping", "application/json", sendData)
			if err != nil {
				if !loopn.PingRetry(sendData) {
					logger.Logger.Warn("Cannot connect to " + loopn.Ip)
					loopn.Active = false
					loopn.PingCount = 0
					return
				}
			}
			if loopn.PingCount == 0 {
				logger.Logger.Info("Sending map to first joiner.")
				jsonNodeMap, _ := json.Marshal(node.MarshalNodeMap())
				sendData := bytes.NewBuffer(jsonNodeMap)
				_, err = http.Post("http://"+loopn.Ip+"/ping", "application/json", sendData)
			}
			loopn.PingCount++
			logger.Logger.Info("Ping received from " + loopn.Ip)
			if loopn.Active == false {
				loopn.Active = true
			}
		}(n)
	}
	node.PingCount++
}

// Every four seconds, we want to ping each connection
func (node *Node) Ping() {
	for {
		time.Sleep(4 * time.Second)
		if node.Rank == global.FOLLOWER {
			return
		}
		if len(node.NodeMap) == 1 {
			continue
		}
		jsonNodeMap := node.CheckIfDataChanged()
		node.PingEachConnection(jsonNodeMap)
	}
}

// We want to check if the master is still alive
func (node *Node) CheckForNoPingFromMaster() {
	if node.Rank == global.MASTER {
		return
	}
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs := time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 1 {
		return
	}
	logger.Logger.Warn("Slow response first check at " + fmt.Sprintf("%f", timeSinceLastPingAbs) + "s")
	time.Sleep(4100 * time.Millisecond)
	timeSinceLastPingAbs = time.Now().Sub(node.Pinged).Seconds()
	if timeSinceLastPingAbs < 8.2 {
		return
	}
	logger.Logger.Info("No ping from master.")
	if node.IsNextInLine() {
		node.SetToMaster()
	}
}

// Set this node to master status and put all 'replicas' to 'active'
func (node *Node) SetToMaster() {
	node.Rank = global.MASTER
	node.Data = node.NodeMap[0].Data
	node.RecentQueries = node.NodeMap[0].RecentQueries
	node.Rules = node.NodeMap[0].Rules

	waitingTimeStr := strconv.Itoa(int(time.Now().Sub(node.Pinged).Seconds()))
	logger.Logger.Info("I'm the master now. I've copied all of the data from the previous master after waiting for " + waitingTimeStr + "s")
	node.NodeMap = node.NodeMap[1:]
	node.NodeMap[0] = node
	go node.Ping()
}

// reads from data storage puts all docs to memory on load
func (node *Node) ReadFromLocal() {
	files, _ := ioutil.ReadDir("data/")
	if len(files) == 0 {
		return
	}
	for _, file := range files {
		key := file.Name()
		data, _ := ioutil.ReadFile("data/" + key)
		var dataJSON global.JsonData
		json.Unmarshal(data, &dataJSON)
		node.Data.Store(key, (map[string]interface{})(dataJSON))
	}
	logger.Logger.Info("Loaded files into memory.")
}

// We want to write a document to disk
func (node *Node) WriteDocToDisk(key string, value interface{}) error {
	dataToWrite, _ := json.Marshal(value)
	path := "data/"

	_ = os.MkdirAll(path, os.ModePerm)

	err := os.WriteFile("data/"+key, []byte(dataToWrite), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
