package global

import (
	"distributed_servers/logging"
	"time"
)

type Node struct {
	Ip            string
	Pinged        time.Time
	PingCount     int
	Rank          Rank
	Data          map[string]interface{}
	Active        bool
	RecentQueries map[string]time.Time
	Rules         map[string][]Task
}

type Rank string
type JsonData map[string]interface{}
type NodeData map[string]interface{}
type Task struct {
	Fnct func(string)
	Data string
}

const (
	MASTER   Rank = "MASTER"
	FOLLOWER Rank = "FOLLOWER"
)

// global vars will clean up
var PORT_START = 3100
var NODE_MAP []*Node
var Logger logging.Logger
var DataChanged bool = false
