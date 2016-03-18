package Session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"io"
)

type RandGenerator struct{}

func (RandGenerator) Generate(name string) string {
	b := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

type GobSerializer struct{}

func (serializer GobSerializer) Unserialize(dst interface{}, reader io.Reader) error {
	return gob.NewDecoder(reader).Decode(dst)
}

func (serializer GobSerializer) Serialize(src interface{}, writer io.Writer) error {
	return gob.NewEncoder(writer).Encode(src)
}
