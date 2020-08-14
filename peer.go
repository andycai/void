package void

type Peer interface {
	PeerProperty
	Start() Peer
	Stop()
	TypeName() string
}

type PeerProperty interface {
	Name() string
	Address() string
	Queue() EventQueue
	SetName(string)
	SetAddress(string)
	SetQueue(EventQueue)
}

type PeerListener interface {
	Port() int
}
