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
