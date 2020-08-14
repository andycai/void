package codec

import (
	"github.com/andycai/void"
	"github.com/andycai/void/meta"
)

func Get(name string) void.Codec {
	for _, c := range codecs {
		if c.Name() == name {
			return c
		}
	}

	return nil
}

func Must(name string) void.Codec {
	c := Get(name)

	if c == nil {
		panic("codec not found " + name)
	}

	return c
}

func Marshal(msg interface{}) (data []byte, m void.Meta, err error) {
	// 获取消息元信息
	m = meta.GetMetaByMsg(msg)
	if m == nil {
		return nil, nil, void.NewErrorContext("msg not exists", msg)
	}

	data, err = m.GetCodec().Marshal(msg)

	return
}

func Unmarshal(msgID int, data []byte) (msg interface{}, m void.Meta, err error) {
	// 获取消息元信息
	m = meta.GetMetaByID(msgID)

	// 消息没有注册
	if m == nil {
		return nil, nil, void.NewErrorContext("msg not exists", msgID)
	}

	// 创建消息
	msg = m.NewType()

	// 从字节数组转换为消息
	err = m.GetCodec().Unmarshal(data, msg)

	return
}
