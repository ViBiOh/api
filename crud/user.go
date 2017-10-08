package crud

type user struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[int64]*user{}
	seq   = int64(1)
)

// Handler for Hello request. Should be use with net/http
type Handler struct {
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
