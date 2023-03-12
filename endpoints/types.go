package endpoints

import "blazem/global"

type Node global.Node

type WebNodeMap struct {
	Ip     string `json:"ip"`
	Active bool   `json:"active"`
}

type SendData struct {
	Key  string                 `json:"key"`
	Data map[string]interface{} `json:"data"`
}

type Stats struct {
	Cpu float64 `json:"cpu"`
	Ram float64 `json:"ram"`
}

type SendQueryData struct {
	Docs      []SendData `json:"docs"`
	TimeTaken int64      `json:"timeTaken"`
}

type Folder struct {
	FolderName string `json:"folderName"`
	DocCount   int    `json:"docCount"`
}
