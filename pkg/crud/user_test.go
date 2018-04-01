package crud

import "testing"
import "reflect"

func Test_listUser(t *testing.T) {
	firstUser := &user{ID: 1, Name: `Test 1`}
	secondUser := &user{ID: 2, Name: `Test 2`}
	thirdUser := &user{ID: 3, Name: `Test 3`}

	var cases = []struct {
		intention string
		init      map[uint]*user
		page      uint
		pageSize  uint
		sortFn    func(*user, *user) bool
		want      []*user
	}{
		{
			`should return page 1`,
			map[uint]*user{1: firstUser},
			1,
			20,
			sortByID,
			[]*user{firstUser},
		},
		{
			`should respect given pageSize`,
			map[uint]*user{1: firstUser, 2: secondUser, 3: thirdUser},
			1,
			2,
			sortByID,
			[]*user{firstUser, secondUser},
		},
		{
			`should respect given page and pageSize`,
			map[uint]*user{1: firstUser, 2: secondUser, 3: thirdUser},
			2,
			2,
			sortByID,
			[]*user{thirdUser},
		},
	}

	for _, testCase := range cases {
		users = testCase.init

		if result := listUser(testCase.page, testCase.pageSize, testCase.sortFn); !reflect.DeepEqual(result, testCase.want) {
			t.Errorf("%v\nlistUser(%v, %v) = %v, want %v", testCase.intention, testCase.page, testCase.pageSize, result, testCase.want)
		}
	}
}
func Test_getUser(t *testing.T) {
	testUser := &user{ID: 1, Name: `Test name`}

	var cases = []struct {
		intention string
		init      map[uint]*user
		id        uint
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
			map[uint]*user{1: testUser},
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
		init      map[uint]*user
		initSeq   uint
		name      string
		want      *user
	}{
		{
			`should create user with given string`,
			map[uint]*user{},
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
		init      map[uint]*user
		id        uint
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
			map[uint]*user{1: {ID: 1, Name: `test`}},
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
		init      map[uint]*user
		id        uint
		want      *user
	}{
		{
			`should do nothing if not found`,
			nil,
			1,
			nil,
		},
		{
			`should remove given instance`,
			map[uint]*user{1: testUser},
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
