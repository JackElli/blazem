package node

import (
	"blazem/pkg/domain/global"
	"blazem/pkg/domain/logger"
	blazem_logger "blazem/pkg/domain/logger"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

const PORT_START = 3100

type INode interface {
	ReadFromLocal()
	WriteDocToDisk(key string, value interface{}) error
	PickPort(ip string) error
	GetLocalIp() string
	SetupLogger() (*zap.Logger, error)
	tryListen(ip string)
}

type Node struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          global.Rank
	Data          sync.Map
	Active        bool
	RecentQueries map[string]string //time
	NodeMap       []*Node
}

func NewNode() *Node {
	return &Node{
		Ip:            "",
		Pinged:        time.Now(),
		PingCount:     0,
		Rank:          global.MASTER,
		Data:          sync.Map{},
		Active:        true,
		RecentQueries: map[string]string{},
	}
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

// We want to pick a port (default 3100) but could try 3 more so max 3103
func (node *Node) PickPort(ip string) error {
	connectIp := ""
	for i := 0; i < 3; i++ {
		connectIp = ip + ":" + strconv.Itoa(PORT_START+i)
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

// Returns the IP of this node
func (node *Node) GetLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// setup file for logging
func (node *Node) SetupLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	blazem_logger.Logger = logger
	defer logger.Sync()
	return logger, nil
}
