package crud

import "sort"

type user struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[uint]*user{}
	seq   = uint(1)
)

type usersSorter struct {
	users []*user
	by    func(p1, p2 *user) bool
}

func (s *usersSorter) Len() int {
	return len(s.users)
}

func (s *usersSorter) Swap(i, j int) {
	s.users[i], s.users[j] = s.users[j], s.users[i]
}

func (s *usersSorter) Less(i, j int) bool {
	return s.by(s.users[i], s.users[j])
}

type sortBy func(p1, p2 *user) bool

func (by sortBy) Sort(arr []*user) {
	sort.Sort(&usersSorter{
		users: arr,
		by:    by,
	})
}

func sortByID(o1, o2 *user) bool {
	return o1.ID < o2.ID
}

func listUser(page, pageSize uint, sortFn func(*user, *user) bool) []*user {
	list := make([]*user, 0)
	for _, value := range users {
		list = append(list, value)
	}

	listSize := uint(len(list))
	sortBy(sortFn).Sort(list)

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

func getUser(id uint) *user {
	return users[id]
}

func createUser(name string) *user {
	createdUser := &user{ID: seq, Name: name}
	users[seq] = createdUser

	seq++
	return createdUser
}

func updateUser(id uint, name string) *user {
	foundUser, ok := users[id]

	if ok {
		foundUser.Name = name
	}

	return foundUser
}

func deleteUser(id uint) {
	delete(users, id)
}
