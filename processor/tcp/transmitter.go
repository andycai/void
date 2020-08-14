package tcp

import (
	"io"
	"net"

	"github.com/andycai/void"
	"github.com/andycai/void/util"
)

type TCPMessageTransmitter struct {
}

func (TCPMessageTransmitter) OnRecvMessage(s void.Session) (msg interface{}, err error) {
	reader, ok := s.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := s.Peer().(void.TCPSocketOption)

	if conn, ok := reader.(net.Conn); ok {
		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {

			msg, err = util.Unpack(reader, opt.MaxPacketSize())

		})
	}

	return
}

func (TCPMessageTransmitter) OnSendMessage(s void.Session, msg interface{}) (err error) {
	writer, ok := s.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := s.Peer().(void.TCPSocketOption)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {
		err = util.Pack(writer, msg)
	})

	return
}
