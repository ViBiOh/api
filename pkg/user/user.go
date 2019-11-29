package user

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/crud"
)

// User describe a user
type User struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// SetID define ID
func (a *User) SetID(id uint64) {
	a.ID = id
}

// Service is a raw implementation of User
type Service struct {
	users map[uint64]*User
	mutex sync.RWMutex
}

// New creates a new user service
func New() *Service {
	return &Service{
		users: map[uint64]*User{},
		mutex: sync.RWMutex{},
	}
}

// Unmarsall a User
func (a *Service) Unmarsall(content []byte) (crud.Item, error) {
	var tag User

	if err := json.Unmarshal(content, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

// List users
func (a *Service) List(_ context.Context, page, pageSize uint, _ string, _ bool, _ map[string][]string) ([]crud.Item, uint, error) {
	list := make([]crud.Item, 0)
	listSize := uint(len(a.users))

	var min uint
	if page > 1 {
		min = (page - 1) * pageSize
	}

	max := min + pageSize
	if max > listSize {
		max = listSize
	}

	if min >= listSize {
		return list, listSize, nil
	}

	for _, value := range a.users {
		list = append(list, value)
	}

	return list[min:max], listSize, nil
}

// Get user by ID
func (a *Service) Get(_ context.Context, id uint64) (crud.Item, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	user, ok := a.users[id]
	if !ok {
		return nil, crud.ErrNotFound
	}

	return user, nil
}

// Create user
func (a *Service) Create(_ context.Context, o crud.Item) (crud.Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user := o.(*User)

	createdUser := &User{ID: uint64(time.Now().Unix()), Name: user.Name}
	a.users[createdUser.ID] = createdUser

	return createdUser, nil
}

// Update user
func (a *Service) Update(_ context.Context, o crud.Item) (crud.Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user := o.(*User)
	a.users[user.ID] = user

	return user, nil
}

// Delete user by ID
func (a *Service) Delete(_ context.Context, o crud.Item) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	delete(a.users, o.(*User).ID)

	return nil
}
