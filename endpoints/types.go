package endpoints

import (
	"blazem/global"
)

type Node global.Node

var connectedFromWebUI bool

// var taskFncDecoder = map[string]func(interface{}, interface{}) (interface{}, error){
// 	"query": func(queryVal interface{}, requirePass interface{}) (interface{}, error) {
// 		queryResult, _, _, _ := query.Execute(queryVal.(string), "")
// 		return queryResult, nil
// 	},
// 	"export": func(hostName interface{}, requirePass interface{}) (interface{}, error) {
// 		getHost, ok := hostName.(string)
// 		if !ok {
// 			return "", fmt.Errorf("not a string host")
// 		}
// 		getDocs, ok := requirePass.([]map[string]interface{})
// 		if !ok {
// 			return "", fmt.Errorf("cannot find docs")
// 		}

// 		if strings.Contains(getHost, "couchbase") {
// 			err := addToCouchbase(getHost, getDocs)
// 			if err != nil {
// 				return "", fmt.Errorf("cannot connect to couchbase")
// 			}
// 		}

// 		return "", nil
// 	},
// }
