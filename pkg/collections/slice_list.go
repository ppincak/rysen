package collections

type SliceList struct {
	// list entries
	entries []interface{}
	// current size
	size int
}

func NewSliceList(initialSize int) *SliceList {
	return &SliceList{
		entries: make([]interface{}, 0),
		size:    0,
	}
}

func (list *SliceList) Entries() []interface{} {
	return list.entries
}

func (list *SliceList) EntriesCopy() []interface{} {
	entries := make([]interface{}, list.Size())
	copy(entries, list.entries)
	return entries
}

func (list *SliceList) Add(value interface{}) {
	list.entries = append(list.entries, value)
	list.size++
}

func (list *SliceList) Get(index int) interface{} {
	if list.size <= index || index < 0 {
		return nil
	} else {
		return list.entries[index]
	}
}

func (list *SliceList) Remove(index int) {
	list.entries = DeletePreserving(list.entries, index)
}

func (list *SliceList) Reset() {
	list.entries = make([]interface{}, 0)
	list.size = 0
}

func (list *SliceList) IsEmpty() bool {
	return list.size == 0
}
func (list *SliceList) Size() int {
	return list.size
}
