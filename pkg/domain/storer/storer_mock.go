package storer

import (
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/node"
	"errors"
)

type StoreMock struct {
	Node          *node.Node
	FolderManager *folder.FolderManager
}

func NewStoreMock(node *node.Node) *StoreMock {
	return &StoreMock{
		Node:          node,
		FolderManager: nil,
	}
}

// Load loads the data from key
func (store *StoreMock) Load(key string) (interface{}, error) {
	value, ok := store.Node.Data.Load(key)
	if !ok {
		return nil, errors.New("Doc not found")
	}
	return value, nil
}

// Store stores a value into a key and increments the folder count
// if there is a folder
func (store *StoreMock) Store(key string, folder interface{}, value interface{}) error {
	store.Node.Data.Store(key, value)
	return nil
}

// Delete deletes a value into a key and decrements the folder count
// if there is a folder
func (store *StoreMock) Delete(key string, folder interface{}) error {
	store.Node.Data.Delete(key)
	return nil
}
