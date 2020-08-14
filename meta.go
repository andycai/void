package void

import "reflect"

type Meta interface {
	TypeName() string
	GetType() reflect.Type
	GetID() int
	NewType() interface{}
	GetCodec() Codec
}
