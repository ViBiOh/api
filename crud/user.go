package crud

import "sort"

type user struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[int64]*user{}
	seq   = int64(1)
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

func listUser(page, pageSize int64, sortFn func(*user, *user) bool) []*user {
	list := make([]*user, 0)
	for _, value := range users {
		list = append(list, value)
	}

	listSize := int64(len(list))
	sortBy(sortFn).Sort(list)

	var min int64
	if page > 1 {
		min = (page - 1) * pageSize
	}
	max := min + pageSize
	if max > listSize {
		max = listSize
	}

	return list[min:max]
}

func getUser(id int64) *user {
	return users[id]
}

func createUser(name string) *user {
	createdUser := &user{ID: seq, Name: name}
	users[seq] = createdUser

	seq++
	return createdUser
}

func updateUser(id int64, name string) *user {
	foundUser, ok := users[id]

	if ok {
		foundUser.Name = name
	}

	return foundUser
}

func deleteUser(id int64) {
	delete(users, id)
}
