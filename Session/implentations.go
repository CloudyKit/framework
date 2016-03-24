package Session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"io"
)

type RandGenerator struct{}
type GobSerializer struct{}

func (RandGenerator) Generate(name string) string {
	b := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
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
