package meta

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/andycai/void"
)

type meta struct {
	Codec void.Codec
	Type  reflect.Type
	ID    int
}

func (m *meta) TypeName() string {
	return m.Type.Name()
}

func (m *meta) FullName() string {
	var sb strings.Builder
	sb.WriteString(path.Base(m.Type.PkgPath()))
	sb.WriteString(".")
	sb.WriteString(m.Type.Name())

	return sb.String()
}

func (m *meta) GetType() reflect.Type {
	return m.Type
}

func (m *meta) GetID() int {
	return m.ID
}

func (m *meta) NewType() interface{} {
	return reflect.New(m.Type).Interface()
}

func (m *meta) GetCodec() void.Codec {
	return m.Codec
}

var (
	metasByID   = make(map[int]void.Meta)
	metasByType = make(map[reflect.Type]void.Meta)
)

func Register(m void.Meta) {
	if m.GetID() == 0 {
		panic("ID could not be 0 by " + m.TypeName())
	}

	if _, ok := metasByType[m.GetType()]; ok {
		panic("duplicate type: " + m.TypeName())
	} else {
		metasByType[m.GetType()] = m
	}

	if _, ok := metasByID[m.GetID()]; ok {
		panic(fmt.Sprintf("duplicate id: %d", m.GetID()))
	} else {
		metasByID[m.GetID()] = m
	}
}

func NewMeta(c void.Codec, t reflect.Type, ID int) void.Meta {
	return &meta{
		Codec: c,
		Type:  t,
		ID:    ID,
	}
}

func GetMetaByID(ID int) void.Meta {
	if v, ok := metasByID[ID]; ok {
		return v
	}

	return nil
}

func GetMetaByType(t reflect.Type) void.Meta {
	if t == nil {
		return nil
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v, ok := metasByType[t]; ok {
		return v
	}

	return nil
}

func GetMetaByMsg(msg interface{}) void.Meta {
	if msg == nil {
		return nil
	}

	return GetMetaByType(reflect.TypeOf(msg))
}

func MessageSize(msg interface{}) int {
	if msg == nil {
		return 0
	}

	m := GetMetaByType(reflect.TypeOf(msg))
	if m == nil {
		return 0
	}

	raw, err := m.GetCodec().Marshal(msg)
	if err != nil {
		return 0
	}

	return len(raw)
}

func MessageToID(msg interface{}) int {
	if msg == nil {
		return 0
	}

	m := GetMetaByMsg(msg)
	if m == nil {
		return 0
	}

	return m.GetID()
}

// 直接发送数据时，将*RawPacket作为Send参数
type RawPacket struct {
	MsgData []byte
	MsgID   int
}

func (r *RawPacket) Message() interface{} {
	// 获取消息元信息
	m := GetMetaByID(r.MsgID)

	// 消息没有注册
	if m == nil {
		return struct{}{}
	}

	// 创建消息
	msg := m.NewType()

	// 从字节数组转换为消息
	err := m.GetCodec().Unmarshal(r.MsgData, msg)
	if err != nil {
		return struct{}{}
	}

	return msg
}
