package processor

import (
	"errors"

	"github.com/andycai/void"
)

var notHandled = errors.New("Processor: Transimitter nil")

type processor struct {
	transmit void.MessageTransmitter
	hooker   void.EventHooker
	handler  void.EventHandler
}

func NewProcessor() void.Processor {
	return &processor{}
}

func (p *processor) SetHooker(h void.EventHooker) {
	p.hooker = h
}

func (p *processor) SetTransmitter(t void.MessageTransmitter) {
	p.transmit = t
}

func (p *processor) SetHandler(handler void.EventHandler) {
	p.handler = handler
}

func (p *processor) ReadMessage(s void.Session) (msg interface{}, err error) {
	if p.transmit != nil {
		return p.transmit.OnRecvMessage(s)
	}

	return nil, notHandled
}

func (p *processor) SendMessage(e void.Event) {
	if p.hooker != nil {
		e = p.hooker.OnOutboundEvent(e)
	}

	if p.transmit != nil && e != nil {
		p.transmit.OnSendMessage(e.Session(), e.Message())
	}
}

func (p *processor) ProcessEvent(e void.Event) {
	if p.hooker != nil {
		e = p.hooker.OnInboundEvent(e)
	}

	if p.handler != nil && e != nil {
		p.handler(e)
	}
}
