package tcp

import (
	"net"
	"time"

	"github.com/andycai/void"
)

type tcpSocketOption struct {
	readBufferSize  int
	writeBufferSize int
	noDelay         bool
	maxPacketSize   int

	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewTCPSocketOption() void.TCPSocketOption {
	return &tcpSocketOption{}
}

func (s *tcpSocketOption) SetSocketBuffer(readBufferSize, writeBufferSize int, noDelay bool) {
	s.readBufferSize = readBufferSize
	s.writeBufferSize = writeBufferSize
	s.noDelay = noDelay
}

func (s *tcpSocketOption) SetSocketDeadline(read, write time.Duration) {
	s.readTimeout = read
	s.writeTimeout = write
}

func (s *tcpSocketOption) SetMaxPacketSize(maxSize int) {
	s.maxPacketSize = maxSize
}

func (s *tcpSocketOption) MaxPacketSize() int {
	return s.maxPacketSize
}

func (s *tcpSocketOption) ApplySocketOption(conn net.Conn) {
	if cc, ok := conn.(*net.TCPConn); ok {

		if s.readBufferSize >= 0 {
			cc.SetReadBuffer(s.readBufferSize)
		}

		if s.writeBufferSize >= 0 {
			cc.SetWriteBuffer(s.writeBufferSize)
		}

		cc.SetNoDelay(s.noDelay)
	}

}

func (s *tcpSocketOption) ApplySocketReadTimeout(conn net.Conn, callback func()) {
	if s.readTimeout > 0 {

		// issue: http://blog.sina.com.cn/s/blog_9be3b8f10101lhiq.html
		conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		callback()
		conn.SetReadDeadline(time.Time{})

	} else {
		callback()
	}
}

func (s *tcpSocketOption) ApplySocketWriteTimeout(conn net.Conn, callback func()) {
	if s.writeTimeout > 0 {

		conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
		callback()
		conn.SetWriteDeadline(time.Time{})

	} else {
		callback()
	}
}

func (s *tcpSocketOption) Init() {
	s.readBufferSize = -1
	s.writeBufferSize = -1
}
