package collections

type List interface {
	Add(interface{})

	Get(index int) interface{}

	Remove(index int)

	Size() int
}
