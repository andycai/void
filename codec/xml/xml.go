package xml

import (
	"encoding/xml"

	"github.com/andycai/void/codec"
)

type xmlCodec struct {
}

func (j *xmlCodec) Name() string {
	return "xml"
}

func (j *xmlCodec) MimeType() string {
	return "application/xml"
}

func (j *xmlCodec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (j *xmlCodec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func init() {
	codec.Register(new(xmlCodec))
}
