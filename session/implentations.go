package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"io"
	"time"
)

// IdGenerator this interface represents an id generator
type IdGenerator interface {
	Generate(id, name string) string
}

type Serializer interface {
	Serialize(src interface{}, w io.Writer) error
	Unserialize(dst interface{}, r io.Reader) error
}

type Store interface {
	Reader(name string) (io.ReadCloser, error)
	Writer(name string) (io.WriteCloser, error)
	Remove(name string) error
	Gc(before time.Time)
}

type RandGenerator struct{}
type GobSerializer struct{}

func (RandGenerator) Generate(id, name string) string {
	if id == "" {
		b := make([]byte, 16)
		n, err := io.ReadFull(rand.Reader, b)
		if n != len(b) || err != nil {
			panic(err)
		}
		return base64.URLEncoding.EncodeToString(b)
	}
	return id
}

func (serializer GobSerializer) Unserialize(dst interface{}, reader io.Reader) error {
	return gob.NewDecoder(reader).Decode(dst)
}

func (serializer GobSerializer) Serialize(src interface{}, writer io.Writer) error {
	return gob.NewEncoder(writer).Encode(src)
}
