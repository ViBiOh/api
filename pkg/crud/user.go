package crud

import (
	"fmt"
	"sync"

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

// UserService is a raw implementation of User
type UserService struct {
	users map[string]*User
	mutex sync.RWMutex
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		users: map[string]*User{},
		mutex: sync.RWMutex{},
	}
}

// Empty returns empy user
func (a *UserService) Empty() Item {
	return &User{}
}

// List users
func (a *UserService) List(page, pageSize uint) []Item {
	list := make([]Item, 0)
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

	return list[min:max]
}

// Get user by ID
func (a *UserService) Get(id string) Item {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.users[id]
}

// Create user
func (a *UserService) Create(o Item) (Item, error) {
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
func (a *UserService) Update(id string, o Item) (Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	foundUser, ok := a.users[id]

	if !ok {
		return nil, ErrNotFound
	}

	foundUser.Name = o.(*User).Name
	return foundUser, nil
}

// Delete user by ID
func (a *UserService) Delete(id string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, ok := a.users[id]

	if !ok {
		return ErrNotFound
	}

	delete(a.users, id)

	return nil
}
