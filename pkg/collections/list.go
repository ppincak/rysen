package collections

type Iterator interface {
	HasNext() bool

	Next() interface{}

	Remove()
}

type List interface {
	Add(interface{})

	AddAt(interface{}, int)

	Get(index int) interface{}

	Remove(index int)

	Size() int

	Iterator() *Iterator
}
