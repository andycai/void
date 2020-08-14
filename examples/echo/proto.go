package main

import (
	"fmt"
	"reflect"

	"github.com/andycai/void/meta"

	"github.com/andycai/void/codec"
	_ "github.com/andycai/void/codec/binary"
	"github.com/andycai/void/util"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

// 将消息注册到系统
func init() {
	meta.Register(meta.NewMeta(
		codec.Must("binary"),
		reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		int(util.StringHash("main.TestEchoACK"))))
}
