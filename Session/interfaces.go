package Session

import (
	"io"
	"time"
)

// IdGenerator this interface represents an id generator
type IdGenerator interface {
	Generate(string) string
}

type Serializer interface {
	Serialize(src interface{}, w io.Writer)
	Unserialize(dst interface{}, r io.Reader)
}

type Store interface {
	Reader(name string) io.ReadCloser
	Writer(name string) io.WriteCloser
	Gc(before time.Time)
}
