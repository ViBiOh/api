package crud

type user struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[int64]*user{}
	seq   = int64(1)
)

func listUser(page, pageSize int64) []*user {
	list := make([]*user, 0, pageSize)

	var i int64
	var min int64

	if page > 1 {
		min = page - 1*pageSize
	}
	max := page * pageSize

	for _, value := range users {
		if min <= i && i < max {
			list = append(list, value)
		}

		if i >= max {
			break
		}
	}

	return list
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
