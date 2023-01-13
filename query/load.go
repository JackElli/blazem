package query

import (
	"distributed_servers/global"
)

// loadTable turns JSON input to one
// we can understand
func loadTable(i interface{}) (interface{}, error) {
	// m := make(map[string]interface{})
	// var ok bool = false

	// _, ok = i.(map[string]interface{})
	// if ok {
	// 	return i.(map[string]interface{}), nil
	// }
	// m = i.(map[string]interface{})
	// for !ok {
	// 	for key, _ := range m {
	// 		if _, gok := m[key].(map[string]interface{}); gok {
	// 			m = m[key].(map[string]interface{})
	// 		}
	// 		if _, fok := m[key].([]interface{}); fok {
	// 			return m[key].([]interface{}), nil
	// 		}

	// 	}
	// }
	// m = i
	return i, nil
}

// LoadIntoMemory loads file or API into mem
func LoadIntoMemory(node global.Node) string {
	// //map of data to return
	// var fileType FileType
	// if strings.Contains(filepath, "http://") ||
	// 	strings.Contains(filepath, "https://") {
	// 	fileType = U
	// } else {
	// 	fileType = F
	// }
	// UploadFileName = filepath
	// // this is var for getting data from file
	// var i interface{}
	// var fdata []byte

	// if fileType == F {
	// 	fdata, _ = ioutil.ReadFile(filepath)
	// } else {
	// 	request, _ := http.Get(filepath)
	// 	fdata, _ = ioutil.ReadAll(request.Body)
	// }

	// err := json.Unmarshal(fdata, &i)
	// if err != nil {

	// 	return "no file uploaded"
	// }
	jsonLoad = node.Data
	return ""

}
