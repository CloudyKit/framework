package registry

type LifeCycle string

type TypeDefinition[Type any] struct {
	provider  func() Type
	disposer  func(Type)
	lifeCycle LifeCycle
}

type Definer[Type any] interface {
	Config(d TypeDefinition[Type]) error
}
