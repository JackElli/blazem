package folder

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	"errors"
	"fmt"

	"encoding/json"
)

type IFolderManager interface {
	IncrementCount(node *node.Node, key string) error
	DecrementCount(node *node.Node, key string) error
	changeCount(node *node.Node, key string, amount int) error
}

type FolderManager struct {
	Node *node.Node
}

func NewFolderManager(node *node.Node) *FolderManager {
	return &FolderManager{
		Node: node,
	}
}

// IncrementCount increments count by 1
func (fm *FolderManager) IncrementCount(key string) error {
	err := fm.changeCount(key, +1)
	return err
}

// DecrementCount decrements count by 1
func (fm *FolderManager) DecrementCount(key string) error {
	err := fm.changeCount(key, -1)
	return err
}

// changeCount changes the docCount by an amount given
func (fm *FolderManager) changeCount(key string, amount int) error {
	doc, ok := fm.Node.Data.Load(key)
	if !ok {
		return errors.New("Cannot load doc")
	}

	folder := doc.(map[string]interface{})
	docCountData, ok := folder["docCount"]
	if !ok {
		return errors.New("Does not have docCount attr")
	}

	docCount := getDocCount(docCountData)
	folder["docCount"] = docCount + amount

	fm.Node.Data.Store(key, folder)
	logger.Logger.Debug(fmt.Sprintf("Set folder: %s docCount to %d", key, docCount))
	return nil
}

// getDocCount returns the int value of docCount
func getDocCount(i interface{}) int {
	switch t := i.(type) {
	case float64:
		return int(t)
	case int:
		return t
	}
	return 0
}

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
