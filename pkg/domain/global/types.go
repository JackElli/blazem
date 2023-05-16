package global

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type Rank string
type JsonData map[string]interface{}
type NodeData map[string]interface{}
type Document map[string]interface{}

const (
	MASTER   Rank = "MASTER"
	FOLLOWER Rank = "FOLLOWER"
)

var PORT_START = 3100
var NODE_MAP []*Node
var Logger *zap.Logger
var DataChanged bool = false
var GlobalNode *Node

type Node struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          Rank
	Data          sync.Map
	Active        bool
	RecentQueries map[string]string //time
	Rules         map[string]Rule
}

type TempNode struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          Rank
	Data          map[string]interface{}
	Active        bool
	RecentQueries map[string]string //time
	Rules         map[string]Rule
}

type Task struct {
	Data    string
	Require int
	Type    string
}

type Replicate struct {
	LocalFolder string `json:"localFolder"`
	RemoteIp    string `json:"remoteIp"`
}

type Rule struct {
	Tasks []JSONTask
	Time  string
}

type JSONTask struct {
	Type    string
	Data    string
	Require int
}
