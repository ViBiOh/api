package crud

import "sort"

type userSorter struct {
	users []*user
	by    func(o1, o2 *user) bool
}

func (s *userSorter) Len() int {
	return len(s.users)
}

func (s *userSorter) Swap(i, j int) {
	s.users[i], s.users[j] = s.users[j], s.users[i]
}

func (s *userSorter) Less(i, j int) bool {
	return s.by(s.users[i], s.users[j])
}

type sortBy func(o1, o2 *user) bool

func (by sortBy) Sort(arr []*user) {
	sort.Sort(&userSorter{
		users: arr,
		by:    by,
	})
}

func sortByID(o1, o2 *user) bool {
	return o1.ID < o2.ID
}

func sortByName(o1, o2 *user) bool {
	return o1.Name < o2.Name
}
