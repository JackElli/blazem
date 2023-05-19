package node

import (
	"blazem/pkg/domain/global"
	blazem_logger "blazem/pkg/domain/logger"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Here, we want to set master attributes and add some sample data when we first
// start with Blazem.
func (node *Node) SetNodeMasterAttrs() {
	node.Rank = global.MASTER
	node.Data = sync.Map{}

	testData1 := map[string]interface{}{
		"type":       "folder",
		"key":        "testkey1",
		"folderName": "TestFolder",
		"value":      "hello this is a test",
		"date":       time.Now().Format("2006-01-02T15:04:05"),
	}
	testData2 := map[string]interface{}{
		"type":   "text",
		"key":    "testkey2",
		"folder": "testkey1",
		"value":  "hello this is a test",
		"date":   time.Now().Format("2006-01-02T15:04:05"),
	}

	node.Data.Store("testkey1", testData1)
	node.Data.Store("testkey2", testData2)
}

// We want to pick a port (default 3100) but could try 3 more so max 3103
func (node *Node) PickPort(ip string) error {
	connectIp := ""
	for i := 0; i < 3; i++ {
		connectIp = ip + ":" + strconv.Itoa(global.PORT_START+i)
		node.tryListen(connectIp)
		if node.Ip != "" {
			break
		}
	}
	return nil
}

// We want to listen on a selected port for this IP
func (node *Node) tryListen(ip string) {
	portstr := ip
	if strings.Count(ip, ":") > 1 {
		portstr = strings.Split(ip, ":")[0]
	}
	blazem_logger.Logger.Info("trying on " + portstr)
	l, err := net.Listen("tcp", portstr)
	if err != nil {
		return
	}
	node.Ip = ip
	blazem_logger.Logger.Info("Blazem started up on " + ip)
	http.Serve(l, nil)
}

// Return the index of the node in the nodemap
func (node *Node) IndexOfNodeInNodeMap() int {
	for i, n := range node.NodeMap {
		if n.Ip == node.Ip {
			return i
		}
	}
	return -1
}

// Return the data stored in the nodemap
func (node *Node) GetNodeDatas() []sync.Map {
	var nodedata []sync.Map
	for _, n := range node.NodeMap {
		nodedata = append(nodedata, n.Data)
	}
	return nodedata
}

// Returns the IP of this node
func (node *Node) GetLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// setup file for logging
func (node *Node) SetupLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	blazem_logger.Logger = logger
	defer logger.Sync()
	return nil
}
