package gob

import (
	"bytes"
	"encoding/gob"

	"github.com/andycai/void/codec"
)

type GobCodec struct {
}

func (g *GobCodec) Name() string {
	return "gob"
}

func (g *GobCodec) MimeType() string {
	return "application/gob"
}

func (g *GobCodec) Marshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *GobCodec) Unmarshal(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

func init() {
	codec.Register(new(GobCodec))
}
