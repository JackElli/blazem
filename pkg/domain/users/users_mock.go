package users

import (
	"blazem/pkg/domain/user"
	"errors"
)

type UserStoreMock struct {
	Users map[string]user.User `json:"users"`
}

func NewUserStoreMock() *UserStoreMock {
	return &UserStoreMock{
		Users: map[string]user.User{},
	}
}

func (us *UserStoreMock) List() map[string]user.User {
	return us.Users
}

func (us *UserStoreMock) SetupUsers() error {
	numOfUsers, err := us.LoadUsers()
	if err != nil {
		return err
	}
	if numOfUsers != 0 {
		return nil
	}
	err = us.Insert("user:1", &user.User{
		Id:       "user:1",
		Name:     "Jack Ellis",
		Username: "JackTest",
		Password: "test123",
		Role:     "admin",
	})
	if err != nil {
		return err
	}
	return nil
}

// LoadUsers loads all of the users stored on disk into memory
// returns the number of users and any errors
func (us *UserStoreMock) LoadUsers() (int, error) {
	usersMock := map[string]user.User{
		"testuser1": {
			Id:       "testuser1",
			Name:     "Jack",
			Role:     "basic",
			Username: "jacke",
			Password: "pass",
		},
	}
	us.Users = usersMock
	return len(us.Users), nil
}

// Get returns a user from the store with the Id passed
func (us *UserStoreMock) Get(id string) (*user.User, error) {
	user, userExists := us.Users[id]
	if !userExists {
		return nil, errors.New("User with that Id does not exist")
	}
	return &user, nil
}

// GetByUsername returns the user and error where the username is equal to
// the username stored
func (us *UserStoreMock) GetByUsername(username string) (*user.User, error) {
	for _, u := range us.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("No user found with that username: " + username)
}

// Insert reads the users file, inserts a user and writes back to disk
func (us *UserStoreMock) Insert(id string, user *user.User) error {
	us.Users[id] = *user
	return nil
}

func (us *UserStoreMock) Update(id string, user *user.User) error {
	return nil
}
