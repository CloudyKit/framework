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
	Serialize(src interface{}, w io.Writer) error
	Unserialize(dst interface{}, r io.Reader) error
}

type Store interface {
	Reader(name string) (io.ReadCloser, error)
	Writer(name string) (io.WriteCloser, error)
	Gc(before time.Time) error
}
