package main

import (
	"blazem/global"
	"blazem/logging"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (node *Node) setNodeMasterAttrs() {
	// Here, we want to set master attributes and add some sample data when we first
	// start with Blazem.
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

func (node *Node) pickPort(ip string) {
	// We want to pick a port (default 3100) but could try 3 more so max 3103
	connectIp := ""
	for i := 0; i < 3; i++ {
		connectIp = ip + ":" + strconv.Itoa(global.PORT_START+i)
		node.tryListen(connectIp)
		if node.Ip != "" {
			break
		}
	}
}

func indexOfNodeInNodeMap(node *global.Node) int {
	// Return the index of the node in the nodemap
	for i, n := range global.NODE_MAP {
		if n.Ip == node.Ip {
			return i
		}
	}
	return -1
}

func getNodeDatas() []sync.Map {
	// Return the data stored in the nodemap
	var nodedata []sync.Map
	for _, n := range global.NODE_MAP {
		nodedata = append(nodedata, n.Data)
	}
	return nodedata
}

func getLocalIp() string {
	// Returns the IP of this node
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

func (node *Node) tryListen(ip string) {
	// We want to listen on a selected port for this IP
	portstr := ip
	if strings.Count(ip, ":") > 1 {
		portstr = strings.Split(ip, ":")[0]
	}
	global.Logger.Log("trying on "+portstr, logging.INFO)
	l, err := net.Listen("tcp", portstr)
	if err != nil {
		return
	}
	node.Ip = ip
	global.Logger.Log("Blazem started up on "+ip, logging.INFO)
	http.Serve(l, nil)
}

func setupLogger() {
	//setup file for logging
	logfile := "logging/"
	global.Logger = *logging.LogFile(logfile)
}
