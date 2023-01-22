package global

import (
	"blazem/logging"
	"sync"
	"time"
)

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

type Rule struct {
	Id          string
	Tasks       []Task
	ExecuteTime *time.Time
	Executed    bool
}

type Replicate struct {
	LocalFolder string `json:"localFolder"`
	RemoteIp    string `json:"remoteIp"`
}

type Rank string
type JsonData map[string]interface{}
type NodeData map[string]interface{}

const (
	MASTER   Rank = "MASTER"
	FOLLOWER Rank = "FOLLOWER"
)

// global vars will clean up
var PORT_START = 3100
var NODE_MAP []*Node
var Logger logging.Logger
var DataChanged bool = false
var GlobalNode *Node
