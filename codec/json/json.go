package json

import (
	"encoding/json"

	"github.com/andycai/void/codec"
)

type jsonCodec struct {
}

func (j *jsonCodec) Name() string {
	return "json"
}

func (j *jsonCodec) MimeType() string {
	return "application/json"
}

func (j *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {
	codec.Register(new(jsonCodec))
}
