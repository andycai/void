package void

import (
	"net"
	"time"
)

type TCPAcceptor interface {
	Peer
	PeerListener
	SessionAccessor
	TCPSocketOption
}

type TCPConnector interface {
	Peer
	PeerListener
	TCPSocketOption

	SetReconnectDuration(time.Duration)
	ReconnectDuration() time.Duration
	Session() Session
	SetSessionManager(raw interface{})
}

type TCPSocketOption interface {
	// 收发缓冲大小，默认-1
	SetSocketBuffer(readBufferSize, writeBufferSize int, noDelay bool)

	// 设置最大的封包大小
	SetMaxPacketSize(maxSize int)

	// 设置读写超时，默认0，不超时
	SetSocketDeadline(read, write time.Duration)

	ApplySocketOption(conn net.Conn)

	MaxPacketSize() int
	ApplySocketReadTimeout(conn net.Conn, callback func())
	ApplySocketWriteTimeout(conn net.Conn, callback func())

	Init()
}
