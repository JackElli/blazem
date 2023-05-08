package node

import (
	"blazem/pkg/domain/global"
	"bytes"
	"sync"
	"time"
)

type INode interface {
	CheckForNoPingFromMaster()
	Ping()
	PingEachConnection(jsonNodeMap []byte)
	PingRetry(sendData *bytes.Buffer) bool
	ReadFromLocal()
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
