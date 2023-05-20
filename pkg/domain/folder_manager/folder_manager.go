package folder_manager

import (
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	"errors"
	"fmt"
)

type IFolderManager interface {
	IncrementCount(node *node.Node, key string) error
	DecrementCount(node *node.Node, key string) error
	changeCount(node *node.Node, key string, amount int) error
	getDocCount(i interface{}) int
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

	docCount := fm.getDocCount(docCountData)
	folder["docCount"] = docCount + amount

	fm.Node.Data.Store(key, folder)
	logger.Logger.Debug(fmt.Sprintf("Set folder: %s docCount to %d", key, docCount))
	return nil
}

// getDocCount returns the int value of docCount
func (fm *FolderManager) getDocCount(i interface{}) int {
	switch t := i.(type) {
	case float64:
		return int(t)
	case int:
		return t
	}
	return 0
}
