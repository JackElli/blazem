package query

import (
	"distributed_servers/global"
)

// loadTable turns JSON input to one
// we can understand
func loadTable(i interface{}) (interface{}, error) {
	return i, nil
}

// LoadIntoMemory loads file or API into mem
func LoadIntoMemory(node global.Node) string {
	jsonLoad = node.Data
	return ""

}
