package crud

import "testing"
import "reflect"

func Test_getUser(t *testing.T) {
	testUser := &user{ID: 1, Name: `Test name`}

	var cases = []struct {
		intention string
		init      map[int64]*user
		id        int64
		want      *user
	}{
		{
			`should work when not found`,
			nil,
			1,
			nil,
		},
		{
			`should return found instance`,
			map[int64]*user{1: testUser},
			1,
			testUser,
		},
	}

	for _, testCase := range cases {
		users = testCase.init

		if result := getUser(testCase.id); result != testCase.want {
			t.Errorf("%v\ngetUser(%v) = %v, want %v", testCase.intention, testCase.id, result, testCase.want)
		}
	}
}

func Test_createUser(t *testing.T) {
	var cases = []struct {
		intention string
		init      map[int64]*user
		initSeq   int64
		name      string
		want      *user
	}{
		{
			`should create user with given string`,
			map[int64]*user{},
			1,
			`test`,
			&user{ID: 1, Name: `test`},
		},
	}

	for _, testCase := range cases {
		users = testCase.init
		seq = testCase.initSeq

		if result := createUser(testCase.name); !reflect.DeepEqual(result, testCase.want) {
			t.Errorf("%v\ncreateUser(%v) = %v, want %v", testCase.intention, testCase.name, result, testCase.want)
		}

		if users[testCase.want.ID] == nil {
			t.Errorf("%v\nNot stored in users", testCase.intention)
		}
	}
}

func Test_updateUser(t *testing.T) {
	var cases = []struct {
		intention string
		init      map[int64]*user
		id        int64
		name      string
		want      *user
	}{
		{
			`should do nothing if not found`,
			nil,
			1,
			`test`,
			nil,
		},
		{
			`should return updated instance`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
			1,
			`Edited test`,
			&user{ID: 1, Name: `Edited test`},
		},
	}

	for _, testCase := range cases {
		users = testCase.init

		if result := updateUser(testCase.id, testCase.name); !reflect.DeepEqual(result, testCase.want) {
			t.Errorf("%v\nupdateUser(%v, %v) = %v, want %v", testCase.intention, testCase.id, testCase.name, result, testCase.want)
		}
	}
}
func Test_deleteUser(t *testing.T) {
	testUser := &user{ID: 1, Name: `Test name`}

	var cases = []struct {
		intention string
		init      map[int64]*user
		id        int64
		want      *user
	}{
		{
			`should do nothing if not found`,
			nil,
			1,
			nil,
		},
		{
			`should properly remove given instance`,
			map[int64]*user{1: testUser},
			1,
			nil,
		},
	}

	for _, testCase := range cases {
		users = testCase.init

		deleteUser(testCase.id)
		if result := users[testCase.id]; result != testCase.want {
			t.Errorf("%v\ndeleteUser(%v) = %v, want %v", testCase.intention, testCase.id, result, testCase.want)
		}
	}
}
