package crud

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ViBiOh/httputils/pkg/uuid"
)

var (
	// ErrUserNotFound occurs when user with given ID if not found
	ErrUserNotFound = errors.New(`User not found`)
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[string]*user{}
	mutex = sync.RWMutex{}
)

func listUser(page, pageSize uint, sortCriteria sortBy) []*user {
	list := make([]*user, 0)
	for _, value := range users {
		list = append(list, value)
	}

	listSize := uint(len(list))
	sortCriteria.Sort(list)

	var min uint
	if page > 1 {
		min = (page - 1) * pageSize
	}
	max := min + pageSize
	if max > listSize {
		max = listSize
	}

	return list[min:max]
}

func getUser(id string) *user {
	mutex.RLock()
	defer mutex.RUnlock()

	return users[id]
}

func createUser(name string) (*user, error) {
	mutex.Lock()
	defer mutex.Unlock()

	newID, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf(`Error while generating UUID: %v`, err)
	}

	createdUser := &user{ID: newID, Name: name}
	users[createdUser.ID] = createdUser

	return createdUser, nil
}

func updateUser(id string, name string) (*user, error) {
	mutex.Lock()
	defer mutex.Unlock()

	foundUser, ok := users[id]

	if !ok {
		return nil, ErrUserNotFound
	}

	foundUser.Name = name
	return foundUser, nil
}

func deleteUser(id string) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := users[id]

	if !ok {
		return ErrUserNotFound
	}

	delete(users, id)

	return nil
}
