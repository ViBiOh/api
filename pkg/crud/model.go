package crud

import "sort"

// Item describe item
type Item interface {
	ID() string
}

// ItemService retrieves item
type ItemService interface {
	Empty() Item
	List(page, pageSize uint) []Item
	Get(ID string) Item
	Create(o Item) (Item, error)
	Update(ID string, o Item) (Item, error)
	Delete(ID string) error
}

// Sorter is a sort implementation
type Sorter struct {
	items []Item
	by    func(o1, o2 Item) bool
}

func (s *Sorter) Len() int {
	return len(s.items)
}

func (s *Sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *Sorter) Less(i, j int) bool {
	return s.by(s.items[i], s.items[j])
}

// SortBy describes comparison function between two items
type SortBy func(o1, o2 Item) bool

// Sort implement the Sort interface
func (s SortBy) Sort(arr []Item) {
	sort.Sort(&Sorter{
		items: arr,
		by:    s,
	})
}
