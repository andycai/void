package binary

import (
	"github.com/andycai/void/codec"
	"github.com/andycai/void/util/binary"
)

type BinaryCodec struct {
}

func (b *BinaryCodec) Name() string {
	return "binary"
}

func (b *BinaryCodec) MimeType() string {
	return "application/binary"
}

func (b *BinaryCodec) Marshal(v interface{}) ([]byte, error) {
	return binary.Marshal(v)
}

func (b *BinaryCodec) Unmarshal(data []byte, v interface{}) error {
	return binary.Unmarshal(data, v)
}

func init() {
	codec.Register(new(BinaryCodec))
}
