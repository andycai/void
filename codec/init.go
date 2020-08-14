package codec

import "github.com/andycai/void"

var codecs []void.Codec

func Register(c void.Codec) {
	if Get(c.Name()) != nil {
		panic("duplicate codec " + c.Name())
	}
	codecs = append(codecs, c)
}
