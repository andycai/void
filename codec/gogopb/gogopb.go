package gogopb

import (
	"errors"

	"github.com/andycai/void/codec"
	"github.com/gogo/protobuf/proto"
)

type GogopbCodec struct {
}

func (g *GogopbCodec) Name() string {
	return "gogopb"
}

func (g *GogopbCodec) MimeType() string {
	return "application/x-protobuf"
}

func (g *GogopbCodec) Marshal(v interface{}) ([]byte, error) {
	if b, ok := v.(proto.Message); ok {
		return proto.Marshal(b)
	}
	return nil, errors.New("v must be proto.Message.")
}

func (g *GogopbCodec) Unmarshal(data []byte, v interface{}) error {
	if o, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, o)
	}
	return errors.New("v must be proto.Message.")
}

func init() {
	codec.Register(new(GogopbCodec))
}
