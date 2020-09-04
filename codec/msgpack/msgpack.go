package msgpack

import (
	"github.com/andycai/void/codec"
	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackCodec struct {
}

func (m *MsgpackCodec) Name() string {
	return "msgpack"
}

func (m *MsgpackCodec) MimeType() string {
	return "application/msgpack"
}

func (m *MsgpackCodec) Marshal(v interface{}) ([]byte, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (m *MsgpackCodec) Unmarshal(data []byte, v interface{}) error {
	if err := msgpack.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func init() {
	codec.Register(new(MsgpackCodec))
}
