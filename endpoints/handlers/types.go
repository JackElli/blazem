package handlers

import (
	"blazem/global"
)

type Node global.Node

type WebNodeMap struct {
	Ip     string `json:"ip"`
	Active bool   `json:"active"`
}

type Stats struct {
	Cpu float64 `json:"cpu"`
	Ram float64 `json:"ram"`
}

type SendData struct {
	Key  string                 `json:"key"`
	Data map[string]interface{} `json:"data"`
}

type SendQueryData struct {
	Docs      []SendData `json:"docs"`
	TimeTaken int64      `json:"timeTaken"`
}

type Folder struct {
	Folder     string  `json:"folder"`
	Key        string  `json:"key"`
	FolderName string  `json:"folderName"`
	DocCount   float64 `json:"docCount"`
	BackedUp   bool    `json:"backedUp"`
}

type JSONTask struct {
	Type    string
	Data    string
	Require int
}

type Rule struct {
	Tasks []JSONTask
	Time  string
}

type DataInFolder struct {
	FolderName    string     `json:"folderName"`
	ParentFolders []Folder   `json:"parentFolders"`
	Data          []SendData `json:"data"`
}

type EndpointResponse struct {
	Code int
	Msg  string
	Data interface{}
}
