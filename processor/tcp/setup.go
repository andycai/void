package tcp

import (
	"github.com/andycai/void"
	"github.com/andycai/void/processor"
)

func init() {
	processor.Register("tcp.ltv", func(p void.Processor, handler void.EventHandler, args ...interface{}) {
		p.SetHooker(new(MsgHooker))
		p.SetTransmitter(new(TCPMessageTransmitter))
		p.SetHandler(processor.NewQueuedEventCallback(handler))
	})
}
