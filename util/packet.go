package util

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/andycai/void/meta"

	"github.com/andycai/void"
	"github.com/andycai/void/codec"
)

var (
	ErrMaxPacket  = errors.New("packet over size")
	ErrMinPacket  = errors.New("packet short size")
	ErrShortMsgID = errors.New("short msgid")
)

const (
	bodySize  = 2 // 包体大小字段
	msgIDSize = 2 // 消息ID字段
)

// 拆包（LTV, Length-Type-Value）
func Unpack(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {
	// Body Size
	var sizeBuffer = make([]byte, bodySize)
	_, err = io.ReadFull(reader, sizeBuffer)

	if err != nil {
		return
	}

	if len(sizeBuffer) < bodySize {
		return nil, ErrMinPacket
	}

	size := binary.LittleEndian.Uint16(sizeBuffer)

	if maxPacketSize > 0 && size >= uint16(maxPacketSize) {
		return nil, ErrMaxPacket
	}

	body := make([]byte, size)

	_, err = io.ReadFull(reader, body)

	if err != nil {
		return
	}

	if len(body) < msgIDSize {
		return nil, ErrShortMsgID
	}

	msgID := binary.LittleEndian.Uint16(body)
	msgData := body[msgIDSize:]

	msg, _, err = codec.Unmarshal(int(msgID), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}

// 封包（LTV, Length-Type-Value）
func Pack(writer io.Writer, data interface{}) error {
	var (
		msgData []byte
		msgID   int
		mt      void.Meta
	)

	switch m := data.(type) {
	case *meta.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error
		msgData, mt, err = codec.Marshal(data)
		if err != nil {
			return err
		}

		msgID = mt.GetID()
	}

	pkt := make([]byte, bodySize+msgIDSize+len(msgData))

	// Length
	binary.LittleEndian.PutUint16(pkt, uint16(msgIDSize+len(msgData)))

	// Type
	binary.LittleEndian.PutUint16(pkt[bodySize:], uint16(msgID))

	// Value (size = bodySize + msgIDSize)
	copy(pkt[bodySize+msgIDSize:], msgData)

	err := WriteFull(writer, pkt)

	return err
}
