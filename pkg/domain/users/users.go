package users

import (
	"blazem/pkg/domain/global"
	"os/user"
)

type UserStorer interface {
	Get(id string) (*user.User, error)
	Insert(id string, user *user.User) error
	Update(id string, user *user.User) error
}

type UserStore struct {
	Node *global.Node
}

func NewUserStore(node *global.Node) *UserStore {
	return &UserStore{
		Node: node,
	}
}

func (u *UserStore) Get(id string) (*user.User, error) {
	return nil, nil
}

func (u *UserStore) Insert(id string, user *user.User) error {
	return nil
}

func (u *UserStore) Update(id string, user *user.User) error {
	return nil
}
