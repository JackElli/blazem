package global

type Rank string

const (
	MASTER   Rank = "MASTER"
	FOLLOWER Rank = "FOLLOWER"
)

type JsonData map[string]interface{}
type NodeData map[string]interface{}
type Document map[string]interface{}

var PORT_START = 3100
var DataChanged bool = false

type Task struct {
	Data    string
	Require int
	Type    string
}

type Replicate struct {
	LocalFolder string `json:"localFolder"`
	RemoteIp    string `json:"remoteIp"`
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
