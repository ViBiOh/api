package user

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/uuid"
	"github.com/pkg/errors"
)

// User describe a user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SetID define ID
func (a *User) SetID(id string) {
	a.ID = id
}

// Service is a raw implementation of User
type Service struct {
	users map[string]*User
	mutex sync.RWMutex
}

// New creates a new user service
func New() *Service {
	return &Service{
		users: map[string]*User{},
		mutex: sync.RWMutex{},
	}
}

// Unmarsall a User
func (a *Service) Unmarsall(content []byte) (crud.Item, error) {
	var tag User

	if err := json.Unmarshal(content, &tag); err != nil {
		return nil, errors.WithStack(err)
	}

	return &tag, nil
}

// List users
func (a *Service) List(_ context.Context, page, pageSize uint, _ string, _ bool, _ map[string][]string) ([]crud.Item, uint, error) {
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

	return list[min:max], listSize, nil
}

// Get user by ID
func (a *Service) Get(_ context.Context, id string) (crud.Item, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.users[id], nil
}

// Create user
func (a *Service) Create(_ context.Context, o crud.Item) (crud.Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user := o.(*User)

	newID, err := uuid.New()
	if err != nil {
		return nil, err
	}

	createdUser := &User{ID: newID, Name: user.Name}
	a.users[createdUser.ID] = createdUser

	return createdUser, nil
}

// Update user
func (a *Service) Update(_ context.Context, o crud.Item) (crud.Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user := o.(*User)

	foundUser, ok := a.users[user.ID]

	if !ok {
		return nil, crud.ErrNotFound
	}

	foundUser.Name = user.Name
	return foundUser, nil
}

// Delete user by ID
func (a *Service) Delete(_ context.Context, o crud.Item) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user := o.(*User)

	delete(a.users, user.ID)

	return nil
}
