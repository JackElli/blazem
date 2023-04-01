package endpoints

import (
	blazem_testing "blazem/blazem_testing/endpoints"
	"reflect"
	"testing"
)

func TestNodeMapHandler(t *testing.T) {
	nodeMapResp := []WebNodeMap{}
	for _, n := range blazem_testing.MOCK_NODE_MAP {
		nodeMapResp = append(nodeMapResp, WebNodeMap{n.Ip, n.Active})
	}

	//testing
	for _, test := range nodeMapResp {
		if reflect.TypeOf(test.Active).String() != "bool" ||
			reflect.TypeOf(test.Ip).String() != "string" {
			t.Error("WebNodeMap not parsed correctly")
		}
	}
}
