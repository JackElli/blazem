package endpoints

import (
	"blazem/global"
	"blazem/query"
	"fmt"
	"strings"
)

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
	Folder     string `json:"folder"`
	Key        string `json:"key"`
	FolderName string `json:"folderName"`
	DocCount   int    `json:"docCount"`
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

type DataInFolder struct {
	ParentFolders []string   `json:"parentFolders"`
	Data          []SendData `json:"data"`
}

var taskFncDecoder = map[string]func(interface{}, interface{}) (interface{}, error){
	"query": func(queryVal interface{}, requirePass interface{}) (interface{}, error) {
		queryResult, _, _, _ := query.Execute(queryVal.(string), "")
		return queryResult, nil
	},
	"export": func(hostName interface{}, requirePass interface{}) (interface{}, error) {
		getHost, ok := hostName.(string)
		if !ok {
			return "", fmt.Errorf("not a string host")
		}
		getDocs, ok := requirePass.([]map[string]interface{})
		if !ok {
			return "", fmt.Errorf("cannot find docs")
		}

		if strings.Contains(getHost, "couchbase") {
			err := addToCouchbase(getHost, getDocs)
			if err != nil {
				return "", fmt.Errorf("cannot connect to couchbase")
			}
		}

		return "", nil
	},
}
