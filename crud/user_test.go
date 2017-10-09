package crud

import "testing"

func TestGetUser(t *testing.T) {
	testUser := &user{ID: 1, Name: `Test name`}

	var cases = []struct {
		init map[int64]*user
		id   int64
		want *user
	}{
		{
			nil,
			1,
			nil,
		},
		{
			map[int64]*user{1: testUser},
			1,
			testUser,
		},
	}

	for _, testCase := range cases {
		users = testCase.init

		if result := getUser(testCase.id); result != testCase.want {
			t.Errorf(`getUser(%v) = %v, want %v`, testCase.id, result, testCase.want)
		}
	}
}
func TestDeleteUser(t *testing.T) {
	testUser := &user{ID: 1, Name: `Test name`}

	var cases = []struct {
		init map[int64]*user
		id   int64
		want *user
	}{
		{
			nil,
			1,
			nil,
		},
		{
			map[int64]*user{1: testUser},
			1,
			nil,
		},
	}

	for _, testCase := range cases {
		users = testCase.init

		deleteUser(testCase.id)
		if result := users[testCase.id]; result != testCase.want {
			t.Errorf(`deleteUser(%v) = %v, want %v`, testCase.id, result, testCase.want)
		}
	}
}
