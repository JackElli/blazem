package storer

import (
	"blazem/pkg/domain/folder"
	"blazem/pkg/domain/node"
	"errors"
)

type Storer interface {
	Load(key string) (interface{}, error)
	Store(key string, folder interface{}, value interface{}) error
	Delete(key string, folder interface{}) error
}

type Store struct {
	Node          *node.Node
	FolderManager *folder.FolderManager
}

func NewStore(node *node.Node) *Store {
	fm := folder.NewFolderManager(node)
	return &Store{
		Node:          node,
		FolderManager: fm,
	}
}

// Load loads the data from key
func (store *Store) Load(key string) (interface{}, error) {
	value, ok := store.Node.Data.Load(key)
	if !ok {
		return nil, errors.New("Doc not found")
	}
	return value, nil
}

// Store stores a value into a key and increments the folder count
// if there is a folder
func (store *Store) Store(key string, folder interface{}, value interface{}) error {
	store.Node.Data.Store(key, value)
	store.Node.WriteDocToDisk(key, value)
	if folder != "" {
		err := store.FolderManager.IncrementCount(folder.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete deletes a value into a key and decrements the folder count
// if there is a folder
func (store *Store) Delete(key string, folder interface{}) error {
	store.Node.Data.Delete(key)
	if folder != "" {
		err := store.FolderManager.DecrementCount(folder.(string))
		if err != nil {
			return err
		}
	}
	return nil
}
