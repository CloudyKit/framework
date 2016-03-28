package common

type Named interface {
	Name() string
}

type URLer interface {
	URL(resource string, v ...interface{}) string
}
