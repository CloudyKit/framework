package meta

type Constructor[T comparable] struct {
	builder func(*T)
}

func (constructor Constructor[T]) New() *T {
	var newObject T
	constructor.builder(&newObject)
	return &newObject
}

func NewConstructor[T comparable](builder func(*T)) Constructor[T] {
	return Constructor[T]{
		builder: builder,
	}
}
