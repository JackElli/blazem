package folder_manager

import (
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	"errors"
	"fmt"
)

// IncrementCount increments count by 1
func IncrementCount(node *node.Node, key string) error {
	err := changeCount(node, key, +1)
	return err
}

// DecrementCount decrements count by 1
func DecrementCount(node *node.Node, key string) error {
	err := changeCount(node, key, -1)
	return err
}

// changeCount changes the docCount by an amount given
func changeCount(node *node.Node, key string, amount int) error {
	doc, ok := node.Data.Load(key)
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

	node.Data.Store(key, folder)
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
