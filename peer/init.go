package peer

import "github.com/andycai/void"

const (
	PEER_TYPE_TCP_ACCEPTOR   = "tcp.Acceptor"
	PEER_TYPE_TCP_CONNECTOR  = "tcp.Connector"
	PEER_TYPE_UDP_ACCEPTOR   = "udp.Acceptor"
	PEER_TYPE_UDP_CONNECTOR  = "udp.Connector"
	PEER_TYPE_WS_ACCEPTOR    = "ws.Acceptor"
	PEER_TYPE_WS_CONNECTOR   = "ws.Connector"
	PEER_TYPE_HTTP_ACCEPTOR  = "http.Acceptor"
	PEER_TYPE_HTTP_CONNECTOR = "http.Connector"
)

// key type name
var peerMap = map[string]func() void.Peer{}

func Register(peerType string, f func() void.Peer) {
	if _, ok := peerMap[peerType]; ok {
		panic("duplicate peer type: " + peerType)
	}

	peerMap[peerType] = f
}

func NewPeer(peerType string, name string, addr string, q void.EventQueue) void.Peer {
	if f, ok := peerMap[peerType]; ok {
		p := f()
		p.SetName(name)
		p.SetAddress(addr)
		p.SetQueue(q)
		return p
	}
	panic("peer type not found [" + peerType + "]")
}
