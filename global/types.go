package global

import (
	"distributed_servers/logging"
	"time"
)

type Node struct {
	Ip        string
	Pinged    time.Time
	PingCount int
	Rank      Rank
	Data      map[string]interface{}
	Active    bool
}

//	type JsonData struct {
//		Key    string    `json:"key"`
//		Folder string    `json:"folder"`
//		Data   string    `json:"data"`
//		Type   string    `json:"type"`
//		Date   time.Time `json:"date"`
//	}
type JsonData map[string]interface{}
type Rank string
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
