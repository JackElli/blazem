package blazem_testing

import (
	"blazem/global"
	"sync"
	"time"
)

var MOCK_NODE_MAP []global.Node = []global.Node{
	{
		Ip:            "192.168.1.1",
		Pinged:        time.Now(),
		PingCount:     1,
		Rank:          global.MASTER,
		Data:          sync.Map{},
		Active:        true,
		RecentQueries: nil,
		Rules:         nil,
	},
}
