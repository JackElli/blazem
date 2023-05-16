package endpoint

import (
	"sync"
)

type EndpointResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// We want to get the length of a sync map
func LenOfSyncMap(mp sync.Map) int {
	var i int
	mp.Range(func(key any, value any) bool {
		i++
		return true
	})
	return i
}
