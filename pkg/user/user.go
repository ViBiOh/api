package user

import (
	"fmt"
	"sync"

	"github.com/ViBiOh/go-api/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/uuid"
)

// User describe a user
type User struct {
	UUID string `json:"id"`
	Name string `json:"name"`
}

// ID returns ID
func (a User) ID() string {
	return a.UUID
}

// Service is a raw implementation of User
type Service struct {
	users map[string]*User
	mutex sync.RWMutex
}

// NewService creates a new user service
func NewService() *Service {
	return &Service{
		users: map[string]*User{},
		mutex: sync.RWMutex{},
	}
}

// Empty returns empy user
func (a *Service) Empty() crud.Item {
	return &User{}
}

// List users
func (a *Service) List(page, pageSize uint) ([]crud.Item, error) {
	list := make([]crud.Item, 0)
	for _, value := range a.users {
		list = append(list, value)
	}

	listSize := uint(len(list))

	var min uint
	if page > 1 {
		min = (page - 1) * pageSize
	}
	max := min + pageSize
	if max > listSize {
		max = listSize
	}

	return list[min:max], nil
}

// Get user by ID
func (a *Service) Get(id string) (crud.Item, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.users[id], nil
}

// Create user
func (a *Service) Create(o crud.Item) (crud.Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user := o.(*User)

	newID, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf(`Error while generating UUID: %v`, err)
	}

	createdUser := &User{UUID: newID, Name: user.Name}
	a.users[createdUser.UUID] = createdUser

	return createdUser, nil
}

// Update user
func (a *Service) Update(id string, o crud.Item) (crud.Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	foundUser, ok := a.users[id]

	if !ok {
		return nil, crud.ErrNotFound
	}

	foundUser.Name = o.(*User).Name
	return foundUser, nil
}

// Delete user by ID
func (a *Service) Delete(id string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, ok := a.users[id]

	if !ok {
		return crud.ErrNotFound
	}

	delete(a.users, id)

	return nil
}
