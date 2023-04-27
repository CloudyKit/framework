package object

import "github.com/CloudyKit/framework/object/meta"

func Constructor[T comparable](builder func(*T)) meta.Constructor[T] {
	return meta.NewConstructor[T](builder)
}
func New[T comparable]() T {
	var t T

	//todo get constructor
	return t
}
