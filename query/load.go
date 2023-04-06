package query

import "blazem/global"

func loadTable(i interface{}) (interface{}, error) {
	// loadTable turns JSON input to one
	// we can understand
	return i, nil
}

func LoadIntoMemory(node global.Node) string {
	// LoadIntoMemory loads file or API into mem
	jsonLoad = node.Data
	return ""
}
