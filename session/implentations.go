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
	Serialize(src interface{}, w io.Writer)
	Unserialize(dst interface{}, r io.Reader)
}

type Store interface {
	Reader(name string) io.ReadCloser
	Writer(name string) io.WriteCloser
	Touch(name string)  error
	//Rename(name string) error

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

func (serializer GobSerializer) Unserialize(dst interface{}, reader io.Reader) {
	err := gob.NewDecoder(reader).Decode(dst)
	if err != nil {
		panic(err)
	}

}

func (serializer GobSerializer) Serialize(src interface{}, writer io.Writer) {
	err := gob.NewEncoder(writer).Encode(src)
	if err != nil {
		panic(err)
	}
}
