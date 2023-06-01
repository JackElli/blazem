package endpoint

import (
	"blazem/pkg/domain/global"
)

type EndpointResponse struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type WebNodeMap struct {
	Ip     string      `json:"ip"`
	Active bool        `json:"active"`
	Rank   global.Rank `json:"rank"`
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

type DataInFolder struct {
	FolderName string     `json:"folderName"`
	Data       []SendData `json:"data"`
}
