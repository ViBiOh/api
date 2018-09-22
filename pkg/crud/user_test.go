package crud

import "testing"
import "reflect"

func Test_listUser(t *testing.T) {
	firstUser := &user{ID: `1`, Name: `Test 1`}
	secondUser := &user{ID: `2`, Name: `Test 2`}
	thirdUser := &user{ID: `3`, Name: `Test 3`}

	var cases = []struct {
		intention string
		init      map[string]*user
		page      uint
		pageSize  uint
		sortFn    func(*user, *user) bool
		want      []*user
	}{
		{
			`should return page 1`,
			map[string]*user{`1`: firstUser},
			1,
			20,
			sortByID,
			[]*user{firstUser},
		},
		{
			`should respect given pageSize`,
			map[string]*user{`1`: firstUser, `2`: secondUser, `3`: thirdUser},
			1,
			2,
			sortByID,
			[]*user{firstUser, secondUser},
		},
		{
			`should respect given page and pageSize`,
			map[string]*user{`1`: firstUser, `2`: secondUser, `3`: thirdUser},
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
