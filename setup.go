package main

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/endpoints"
	"blazem/pkg/query"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type SetupManager struct {
	Steps []SetupStep
	Node  *Node
}

type SetupStep struct {
	Description string
	Fn          func() error
}

// Returns a setupmgr with the steps to complete and the node
func (node *Node) CreateSetupMgr(steps []SetupStep) SetupManager {
	return SetupManager{
		Steps: steps,
		Node:  node,
	}
}

// Runs all the steps in order
func (mgr *SetupManager) RunSteps() {
	global.Logger.Info("Setting up Blazem.")
	for _, step := range mgr.Steps {
		if err := step.Fn(); err != nil {
			global.Logger.Error("Found error in " + step.Description + " " + err.Error())
			return
		}
		global.Logger.Info("Completed step.")
	}
	global.Logger.Info("All steps completed successfully :)")
}

// Run the setup process by creating a setup mgr and running each
// step
func (node *Node) RunSetup() {
	var masterip string = ""
	var localip = getLocalIp()
	global.GlobalNode = (*global.Node)(node)

	setupLogger()

	mgr := node.CreateSetupMgr([]SetupStep{
		{
			"Picks port for blazem to start on",
			func() error {
				go node.pickPort(localip)
				return nil
			},
		},
		{
			"Sets up blazem endpoints",
			func() error {
				if err := endpoints.SetupEndpoints((*global.Node)(node)); err != nil {
					return err
				}
				return nil
			},
		},
		{
			"Adds this node to the nodemap",
			func() error {
				global.NODE_MAP = append(global.NODE_MAP, (*global.Node)(node))
				return nil
			},
		},
		{
			"If this node is the master, set master attrs",
			func() error {
				if masterip == node.Ip {
					node.setNodeMasterAttrs()
				}
				return nil
			},
		},
		{
			"Read from local storage",
			func() error {
				(*global.Node)(node).ReadFromLocal()
				return nil
			},
		},
		{
			"First ping and ping either the master or followers",
			func() error {
				go (*global.Node)(node).Ping()
				return nil
			},
		},
		{
			"Load all query data into memory",
			func() error {
				query.LoadIntoMemory(global.Node(*node))
				return nil
			},
		},
	})
	mgr.RunSteps()
	// go (*endpoints.Node)(&node).CheckRules()
}

// Here, we want to set master attributes and add some sample data when we first
// start with Blazem.
func (node *Node) setNodeMasterAttrs() {
	node.Rank = global.MASTER
	node.Data = sync.Map{}

	testData1 := global.Document{
		"type":       "folder",
		"key":        "testkey1",
		"folderName": "TestFolder",
		"value":      "hello this is a test",
		"date":       time.Now().Format("2006-01-02T15:04:05"),
	}
	testData2 := global.Document{
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
func (node *Node) pickPort(ip string) error {
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
	global.Logger.Info("trying on " + portstr)
	l, err := net.Listen("tcp", portstr)
	if err != nil {
		return
	}
	node.Ip = ip
	global.Logger.Info("Blazem started up on " + ip)
	http.Serve(l, nil)
}

// Return the index of the node in the nodemap
func indexOfNodeInNodeMap(node *global.Node) int {
	for i, n := range global.NODE_MAP {
		if n.Ip == node.Ip {
			return i
		}
	}
	return -1
}

// Return the data stored in the nodemap
func getNodeDatas() []sync.Map {
	var nodedata []sync.Map
	for _, n := range global.NODE_MAP {
		nodedata = append(nodedata, n.Data)
	}
	return nodedata
}

// Returns the IP of this node
func getLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// setup file for logging
func setupLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	global.Logger = logger
	defer logger.Sync()
	return nil
}
