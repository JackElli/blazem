package global

type Rank string

const (
	MASTER   Rank = "MASTER"
	FOLLOWER Rank = "FOLLOWER"
)

type JsonData map[string]interface{}

var DataChanged bool = false
