package void

type Event interface {
	Session() Session
	Message() interface{}
}

type MessageTransmitter interface {
	OnRecvMessage(s Session) (msg interface{}, err error)
	OnSendMessage(s Session, msg interface{}) error
}

type EventHooker interface {
	OnInboundEvent(input Event) (output Event)
	OnOutboundEvent(input Event) (output Event)
}

type EventHandler func(e Event)

type Processor interface {
	SetHooker(h EventHooker)
	SetTransmitter(t MessageTransmitter)
	SetHandler(handler EventHandler)
	ReadMessage(s Session) (msg interface{}, err error)
	SendMessage(e Event)
	ProcessEvent(e Event)
}
