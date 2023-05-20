package folder

import (
	"blazem/pkg/domain/endpoint"

	"encoding/json"
)

// FolderToMap converts a folder type to map so we can store it
func FolderToMap(folder endpoint.Folder) (map[string]interface{}, error) {
	bytes, err := json.Marshal(folder)
	if err != nil {
		return nil, err
	}
	var _map map[string]interface{}
	err = json.Unmarshal(bytes, &_map)
	if err != nil {
		return nil, err
	}
	_map["type"] = "folder"
	return _map, nil
}

// IsFolder returns a folder and a bool true if it is a folder and false
// if not
func IsFolder(_map map[string]interface{}) (*endpoint.Folder, bool) {
	mapData, err := json.Marshal(_map)
	if err != nil {
		return nil, false
	}
	var folder endpoint.Folder
	err = json.Unmarshal(mapData, &folder)
	if folder.Type != "folder" {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	return &folder, true
}
