package users

import (
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/user"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type UserStorer interface {
	SetupUsers() error
	LoadUsers() (int, error)
	Get(id string) (*user.User, error)
	GetByUsername(username string) (*user.User, error)
	Insert(id string, user *user.User) error
	Update(id string, user *user.User) error
}

type UserStore struct {
	Users map[string]user.User `json:"users"`
}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (us *UserStore) SetupUsers() error {
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
func (us *UserStore) LoadUsers() (int, error) {
	data, err := ioutil.ReadFile("/users/users.json")
	if err != nil {
		err = os.MkdirAll("/users/", os.ModePerm)
		if err != nil {
			return 0, nil
		}
		err = ioutil.WriteFile("/users/users.json", []byte("{\"users\":{}}"), 0x77)
		logger.Logger.Info("Writing new user file.")
		if err != nil {
			return 0, nil
		}
		data, err = ioutil.ReadFile("/users/users.json")
		if err != nil {
			return 0, nil
		}
	}
	var users UserStore
	err = json.Unmarshal(data, &users)
	if err != nil {
		return 0, err
	}
	us.Users = users.Users
	return len(us.Users), err
}

// Get returns a user from the store with the Id passed
func (us *UserStore) Get(id string) (*user.User, error) {
	user, userExists := us.Users[id]
	if !userExists {
		return nil, errors.New("User with that Id does not exist")
	}
	return &user, nil
}

// GetByUsername returns the user and error where the username is equal to
// the username stored
func (us *UserStore) GetByUsername(username string) (*user.User, error) {
	for _, u := range us.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("No user found with that username: " + username)
}

// Insert reads the users file, inserts a user and writes back to disk
func (us *UserStore) Insert(id string, user *user.User) error {
	data, err := ioutil.ReadFile("/users/users.json")
	if err != nil {
		return err
	}

	var users UserStore
	err = json.Unmarshal(data, &users)
	if err != nil {
		return err
	}

	_, userExists := us.Users[id]
	if userExists {
		return errors.New("User already exists with that Id")
	}

	us.Users[id] = *user
	dataToWrite, err := json.Marshal(UserStore{
		Users: us.Users,
	})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/users/users.json", dataToWrite, 0x77)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Update(id string, user *user.User) error {
	return nil
}
