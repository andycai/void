package void

type Session interface {
	Raw() interface{}
	Peer() Peer
	Send(msg interface{})
	Close()
	ID() int64
	SetID(int64)
}

type SessionAccessor interface {
	GetSession(int64) Session
	VisitSession(func(Session) bool)
	SessionCount() int
	CloseAllSession()
}

type SessionManager interface {
	SessionAccessor
	Add(s Session)
	Remove(s Session)
	Count() int
	SetIDBase(base int64)
}
